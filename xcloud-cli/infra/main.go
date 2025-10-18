package main

import (
	"os"

	"github.com/pulumi/pulumi-github/sdk/v6/go/github"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		owner := os.Getenv("GITHUB_OWNER")
		repoName := "XControl"

		prov, err := github.NewProvider(ctx, "github", &github.ProviderArgs{
			Owner: pulumi.StringPtr(owner),
			Token: pulumi.StringPtr(os.Getenv("GITHUB_TOKEN")),
		})
		if err != nil {
			return err
		}

		repo, err := github.NewRepository(ctx, "managed-repo", &github.RepositoryArgs{
			Name:                pulumi.String(repoName),
			AllowMergeCommit:    pulumi.Bool(false),
			AllowRebaseMerge:    pulumi.Bool(true),
			AllowSquashMerge:    pulumi.Bool(true),
			DeleteBranchOnMerge: pulumi.Bool(true),
		}, pulumi.Provider(prov))
		if err != nil {
			return err
		}

		_, err = github.NewRepositoryRuleset(ctx, "protect-release-pattern", &github.RepositoryRulesetArgs{
			Name:        pulumi.String("Protect release/*"),
			Repository:  repo.Name,
			Target:      pulumi.String("branch"),
			Enforcement: pulumi.String("active"),
			Conditions: &github.RepositoryRulesetConditionsArgs{
				RefName: &github.RepositoryRulesetConditionsRefNameArgs{
					Includes: pulumi.StringArray{
						pulumi.String("release/*"),
					},
					Excludes: pulumi.StringArray{},
				},
			},
			Rules: &github.RepositoryRulesetRulesArgs{
				PullRequest: &github.RepositoryRulesetRulesPullRequestArgs{
					RequiredApprovingReviewCount:   pulumi.Int(1),
					RequireCodeOwnerReview:         pulumi.Bool(false),
					RequireLastPushApproval:        pulumi.Bool(false),
					DismissStaleReviewsOnPush:      pulumi.Bool(true),
					RequiredReviewThreadResolution: pulumi.Bool(true),
				},
				RequiredLinearHistory: pulumi.Bool(true),
				NonFastForward:        pulumi.Bool(true),
				RequiredStatusChecks: &github.RepositoryRulesetRulesRequiredStatusChecksArgs{
					StrictRequiredStatusChecksPolicy: pulumi.Bool(true),
					DoNotEnforceOnCreate:             pulumi.Bool(false),
					RequiredChecks: github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArray{
						&github.RepositoryRulesetRulesRequiredStatusChecksRequiredCheckArgs{
							Context: pulumi.String("require-cherrypick"),
						},
					},
				},
			},
		}, pulumi.Provider(prov))
		if err != nil {
			return err
		}

		return nil
	})
}
