package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

const (
	gmxcliAuthConfig = "auth.json"
)

// These will be over-ridden by the build process. See build.sh.
var (
	clientID     = "invalid"
	clientSecret = "invalid"
)

// TokenMap maps user IDs (typically email) to oauth tokens.
type TokenMap map[string]*oauth2.Token

// GMXConfig captures configuration for gmxcli commands
type GMXConfig struct {
	Filters []gmail.Filter
}

var svc *gmail.Service

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:               "gmxcli",
	SilenceUsage:      true,
	PersistentPreRunE: initGmx,
}

// Execute adds all child commands to the root command sets flags appropriately.
// This is called by gmxcli.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}

func init() {
	RootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging")
	RootCmd.PersistentFlags().StringP("user", "u", "me", "User to authenticate as")
	RootCmd.PersistentFlags().StringP("config", "c", "", "YAML configuration file (required)")
	RootCmd.MarkFlagRequired("config")
}

func initGmx(cmd *cobra.Command, args []string) error {
	user, _ := cmd.Flags().GetString("user")
	client := newOAuthClient(user)
	svc, _ = gmail.New(client)
	if svc == nil {
		log.Fatalf("Unable to create Gmail service")
	}
	return nil
}

// UserCacheDir returns the default root directory to use for user-specific
// cached data. Users should create their own application-specific subdirectory
// within this one and use that.
//
// On Unix systems, it returns $XDG_CACHE_HOME as specified by
// https://standards.freedesktop.org/basedir-spec/basedir-spec-latest.html if
// non-empty, else $HOME/.cache.
// On Darwin, it returns $HOME/Library/Caches.
// On Windows, it returns %LocalAppData%.
// On Plan 9, it returns $home/lib/cache.
//
// If the location cannot be determined (for example, $HOME is not defined),
// then it will return an error.
// TODO(diwaker): remove this in favor of os.UserCacheDir once go 1.11 is released
func UserCacheDir() (string, error) {
	var dir string

	switch runtime.GOOS {
	case "windows":
		dir = os.Getenv("LocalAppData")
		if dir == "" {
			return "", errors.New("%LocalAppData% is not defined")
		}

	case "darwin":
		dir = os.Getenv("HOME")
		if dir == "" {
			return "", errors.New("$HOME is not defined")
		}
		dir += "/Library/Caches"

	case "plan9":
		dir = os.Getenv("home")
		if dir == "" {
			return "", errors.New("$home is not defined")
		}
		dir += "/lib/cache"

	default: // Unix
		dir = os.Getenv("XDG_CACHE_HOME")
		if dir == "" {
			dir = os.Getenv("HOME")
			if dir == "" {
				return "", errors.New("neither $XDG_CACHE_HOME nor $HOME are defined")
			}
			dir += "/.cache"
		}
	}

	return dir, nil
}

func readTokens(filePath string) (TokenMap, error) {
	b, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	var tm TokenMap
	if json.Unmarshal(b, &tm) != nil {
		return nil, err
	}

	return tm, nil
}

func newOAuthClient(user string) *http.Client {
	cacheDir, err := UserCacheDir()
	if err != nil {
		return nil
	}

	filePath := path.Join(cacheDir, "gmxcli", gmxcliAuthConfig)
	tokenMap, _ := readTokens(filePath)
	if tokenMap == nil {
		tokenMap = make(TokenMap)
	}

	ctx := context.Background()
	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{gmail.GmailLabelsScope, gmail.GmailSettingsBasicScope},
	}
	token, found := tokenMap[user]
	if !found {
		token = tokenFromWeb(ctx, config)
		tokenMap[user] = token
		writeTokens(filePath, tokenMap)
	} else {
		log.Printf("Using cached token from %q", filePath)
	}
	return config.Client(ctx, token)
}

func tokenFromWeb(ctx context.Context, config *oauth2.Config) *oauth2.Token {
	ch := make(chan string)
	randState := fmt.Sprintf("st%d", time.Now().UnixNano())
	ts := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/favicon.ico" {
			http.Error(rw, "", 404)
			return
		}
		if req.FormValue("state") != randState {
			log.Printf("State doesn't match: req = %#v", req)
			http.Error(rw, "", 500)
			return
		}
		if code := req.FormValue("code"); code != "" {
			fmt.Fprintf(rw, "<h1>Success</h1>Authorized.")
			rw.(http.Flusher).Flush()
			ch <- code
			return
		}
		log.Printf("no code")
		http.Error(rw, "", 500)
	}))
	defer ts.Close()

	config.RedirectURL = ts.URL
	authURL := config.AuthCodeURL(randState)
	log.Printf("Authorize this app at: %s", authURL)
	code := <-ch
	log.Printf("Got code: %s", code)

	token, err := config.Exchange(ctx, code)
	if err != nil {
		log.Fatalf("Token exchange error: %v", err)
	}
	return token
}

func writeTokens(filePath string, tokens TokenMap) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Doesn't exist; lets create it
		err = os.MkdirAll(filepath.Dir(filePath), 0700)
		if err != nil {
			return
		}
	}

	// At this point, file must exist. Lets (over)write it.
	b, err := json.Marshal(tokens)
	if err != nil {
		return
	}
	if err = ioutil.WriteFile(filePath, b, 0600); err != nil {
		return
	}
}
