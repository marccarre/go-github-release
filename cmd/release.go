package cmd

import (
	"context"
	"os"
	"path/filepath"

	"github.com/marccarre/go-github-release/pkg/gpg"
	"github.com/marccarre/go-github-release/pkg/validate"

	"github.com/google/go-github/github"
	"github.com/marccarre/go-github-release/pkg/logging"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

// releaseCmd represents the release command
var releaseCmd = &cobra.Command{
	Use:   "release",
	Short: "Sign and upload the provided release assets on GitHub under the release corresponding to the provided tag",
	Args:  cobra.MinimumNArgs(1), // At least one release asset to sign and upload.
	Run:   releaseRun,
}

var owner string
var repo string
var tag string
var gpgPrivKey string
var draft bool

func init() {
	zapLogger := logging.MustZapLogger()
	defer zapLogger.Sync() // Flushes buffer, if any.
	log := zapLogger.Sugar()

	releaseCmd.Flags().StringVarP(&owner, "owner", "o", "", "GitHub owner, e.g. marccarre in github.com/marccarre/go-github-release")
	mustMarkFlagRequired(log, "owner")
	releaseCmd.Flags().StringVarP(&repo, "repo", "r", "", "GitHub repository, e.g. go-github-release in github.com/marccarre/go-github-release")
	mustMarkFlagRequired(log, "repo")
	releaseCmd.Flags().StringVarP(&tag, "tag", "t", "", "Git tag corresponding to the release to perform, e.g. v1.0.0")
	mustMarkFlagRequired(log, "tag")
	releaseCmd.Flags().StringVarP(&gpgPrivKey, "key", "k", "", "Path to the private GPG key to use to sign the release assets")
	mustMarkFlagRequired(log, "key")
	releaseCmd.Flags().BoolVarP(&draft, "draft", "d", true, "Should the release be a draft release, default: true")
	rootCmd.AddCommand(releaseCmd)
}

func mustMarkFlagRequired(log *zap.SugaredLogger, flag string) {
	if err := releaseCmd.MarkFlagRequired(flag); err != nil {
		log.Fatalw("failed to mark flag as required", "flag", flag, "error", err)
	}
}

const (
	envGithubToken = "GITHUB_API_TOKEN" // #nosec: environment variable's name. Fixes "Potential hardcoded credentials,HIGH,LOW (gosec)"
	envGPGPassword = "GPG_PASSWD"
)

func releaseRun(cmd *cobra.Command, filePaths []string) {
	zapLogger := logging.MustZapLogger()
	defer zapLogger.Sync() // Flushes buffer, if any.
	log := zapLogger.Sugar()

	mustValidateInputs(log, filePaths)

	ctx := context.Background()
	client := newGitHubClient(ctx, os.Getenv(envGithubToken))
	signer, err := gpg.NewSigner(gpgPrivKey, os.Getenv(envGPGPassword))
	if err != nil {
		log.Fatalw("failed to create GPG signer", "error", err)
	}

	release := mustCreateRelease(ctx, log, client)
	for _, filePath := range filePaths {
		signaturePath := mustSignReleaseAsset(log, signer, filePath)
		mustUploadReleaseAsset(ctx, log, client, filePath, release)
		mustUploadReleaseAsset(ctx, log, client, signaturePath, release)
	}
}

func mustValidateInputs(log *zap.SugaredLogger, filePaths []string) {
	if err := validate.Env(envGithubToken); err != nil {
		log.Fatalw("invalid GitHub API token", "error", err)
	}
	if err := validate.Env(envGPGPassword); err != nil {
		log.Fatalw("invalid GitHub API token", "error", err)
	}
	if err := validate.Files(filePaths); err != nil {
		log.Fatalw("invalid files", "error", err)
	}
	if err := validate.File(gpgPrivKey); err != nil {
		log.Fatalw("invalid GPG private key", "error", err)
	}
}

func newGitHubClient(ctx context.Context, accessToken string) *github.Client {
	tokenSource := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: accessToken},
	)
	oauthClient := oauth2.NewClient(ctx, tokenSource)
	return github.NewClient(oauthClient)
}

func mustCreateRelease(ctx context.Context, log *zap.SugaredLogger, client *github.Client) *github.RepositoryRelease {
	log.Infow("creating release", "owner", owner, "repo", repo, "tag", tag, "draft", draft)
	release, _, err := client.Repositories.CreateRelease(ctx, owner, repo, &github.RepositoryRelease{
		TagName: &tag,
		Name:    &tag,
		Draft:   &draft,
	})
	if err != nil {
		log.Fatalw("failed to create release", "owner", owner, "repo", repo, "tag", tag, "draft", draft, "error", err)
	}
	log.Infow("successfully created release", "owner", owner, "repo", repo, "tag", tag, "draft", draft)
	return release
}

func mustSignReleaseAsset(log *zap.SugaredLogger, signer *gpg.Signer, filePath string) string {
	log.Infow("signing release asset", "file", filePath)
	signaturePath, err := signer.ArmoredDetachSign(filePath)
	if err != nil {
		log.Fatalw("failed to sign release asset", "file", filePath, "error", err)
	}
	log.Infow("successfully signed release asset", "file", filePath)
	return signaturePath
}

func mustUploadReleaseAsset(ctx context.Context, log *zap.SugaredLogger, client *github.Client, filePath string, release *github.RepositoryRelease) {
	fileName := filepath.Base(filePath)
	log.Infow("uploading release asset", "file", fileName, "path", filePath, "release", release.Name)
	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		log.Fatalw("failed to open file", "file", fileName, "path", filePath, "error", err)
	}
	defer file.Close()
	asset, _, err := client.Repositories.UploadReleaseAsset(ctx, owner, repo, *release.ID, &github.UploadOptions{
		Name: fileName,
	}, file)
	if err != nil {
		log.Fatalw("failed to upload release asset", "file", fileName, "path", filePath, "error", err)
	}
	log.Infow("successfully uploaded release asset", "file", fileName, "path", filePath, "release", release.Name, "asset", asset.Name)
}
