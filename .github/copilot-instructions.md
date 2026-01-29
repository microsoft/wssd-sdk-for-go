# GitHub Copilot Instructions for wssd-sdk-for-go

## PR Readiness Requirements

When working on pull requests in this repository, **DO NOT mark the PR as ready for review** until ALL of the following conditions are met:

### 1. Local Build and Test Verification

Before marking a PR ready for review, you MUST verify that all local checks pass:

- **Build**: Run `make build` and ensure it completes successfully with no errors
- **Unit Tests**: Run `make unittest` and ensure all tests pass
- **Linting**: Run `make golangci-lint` and ensure no linting errors are reported
- **Formatting**: Run `make format` to ensure code is properly formatted

### 2. CI/CD Pipeline Status

You MUST verify that all continuous integration checks are passing:

- **Azure Pipelines**: Check that all pipeline jobs have completed successfully
  - Build job (includes `make all`)
  - Unit Tests job
  - Lint job (runs golangci-lint)
  - Static Analysis job (security scans)

- **GitHub Actions**: Check that all GitHub Actions workflows are passing
  - CodeQL analysis
  - Dependency review
  - Any other configured workflows

### 3. How to Check CI Status

Use the GitHub MCP tools to verify status:

1. List recent workflow runs: `list_workflow_runs` 
2. Check specific run status: `get_workflow_run`
3. Review job logs if failures occur: `get_job_logs`
4. For Azure Pipelines, check the PR page for pipeline status

### 4. Required Actions Before Marking Ready

1. **Wait for all checks**: Do not mark the PR as ready until all automated checks show green/passing status
2. **Fix failures**: If any check fails, investigate the logs, fix the issues, commit changes, and wait for checks to run again
3. **Verify build quality**: Ensure no warnings or errors in build output that could indicate issues
4. **Security checks**: Ensure all security scanning tools (CredScan, PoliCheck, GoSec, CodeQL) pass

### 5. General Guidelines

- Always run local validation before pushing to reduce CI failures
- If a check fails repeatedly, investigate root cause rather than repeatedly pushing
- Document any known issues or limitations in the PR description
- Ensure your changes don't break existing tests or introduce new security vulnerabilities

## Build and Test Commands

For reference, here are the key commands for this repository:

```bash
# Format code
make format

# Build the project
make build

# Run unit tests
make unittest

# Run linting
make golangci-lint

# Run all checks (format, build, test)
make all

# Clean build artifacts
make clean
```

## Important Notes

- This repository uses Go 1.25 (as specified in Azure Pipelines)
- Module mode is enabled (GO111MODULE=on)
- Private repo workaround is configured for Microsoft GitHub repos
- The default branch is `master`
- PRs must pass security analysis including CredScan, PoliCheck, GoSec, and CodeQL

---

**Remember**: Green status across ALL checks is mandatory before marking any PR as ready for review. This ensures code quality, security, and prevents breaking changes from being merged.
