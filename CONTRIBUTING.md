# Contributing to Go URL Shortener

First off, thank you for considering contributing to Go URL Shortener. It's people like you that make it such a great tool.

## How Can I Contribute?

### Reporting Bugs

This section guides you through submitting a bug report. Following these guidelines helps maintainers and the community understand your report, reproduce the behavior, and find related reports.

- **Use a clear and descriptive title** for the issue to identify the problem.
- **Provide a step-by-step description** of the issue in as many details as possible.
- **Provide specific examples** to demonstrate the steps. Include links to files or GitHub projects, or copy/pasteable snippets, which you use in those examples.
- **Describe the behavior you observed** after following the steps and point out what exactly is the problem with that behavior.
- **Explain which behavior you expected** to see instead and why.

### Suggesting Enhancements

This section guides you through submitting an enhancement suggestion, including completely new features and minor improvements to existing functionality.

- **Use a clear and descriptive title** for the issue to identify the suggestion.
- **Provide a step-by-step description** of the suggested enhancement in as many details as possible.
- **Provide specific examples** to demonstrate the steps or the expected behavior.
- **Describe the current behavior** and **explain which behavior you expected** to see instead and why.

### Pull Requests

The process described here has several goals:

- Maintain the quality of the project.
- Fix problems that are important to users.
- Engage the community in working toward the best possible version of Go URL Shortener.
- Enable a sustainable system for maintainers to review contributions.

Please follow these steps for your contribution:

1. **Fill in the required template**: All pull requests should include information about the changes, including the purpose of the update and the context of the changes.

2. **Follow the coding conventions**: The code style should be consistent with the rest of the codebase.

3. **Consider the scope of your change**: If your change is large or complex, we may ask you to break it down into smaller pull requests.

4. **Write meaningful commit messages**: Include a brief description of your changes in the commit messages.

5. **Include tests**: Ensure new and existing features run correctly by updating and writing new tests.

### Code Quality Rules (Most Important)

To ensure the codebase remains clean, readable, and maintainable, please adhere to the following rules:

- **Pass all tests**: Your code must pass all existing tests and any new tests you have written to cover your changes.

- **Keep a low cyclomatic complexity**: We aim for a cyclomatic complexity of `5` or less. Use tools like `gocyclo` to check the complexity of your code.

- **Follow Go best practices**: Refer to resources like "Effective Go" and the Go code review comments for guidance on writing idiomatic Go code.

- **Avoid code smells**: Refactor any code that might be considered a "code smell" - a surface indication that usually corresponds to a deeper problem in the system.


## Pull Request Process

1. Ensure any install or build dependencies are removed before the end of the layer when doing a build.
2. Update the README.md with details of changes to the interface, this includes new environment variables, exposed ports, useful file locations, and container parameters.
3. Increase the version numbers in any examples files and the README.md to the new version that this Pull Request would represent. The versioning scheme we use is [SemVer](http://semver.org/).
4. You may merge the Pull Request in once you have the sign-off of two other developers, or if you do not have permission to do that, you may request the second reviewer to merge it for you.

## Styleguides

### Git Commit Messages

- Use the present tense ("Add feature" not "Added feature").
- Use the imperative mood ("Move cursor to..." not "Moves cursor to...").
- Limit the first line to 72 characters or less.
- Reference issues and pull requests liberally after the first line.

### Go Styleguide

- The project follows the standard Go conventions as outlined in [Effective Go](https://go.dev/doc/effective_go) and [Go Code Review Comments](https://go.dev/wiki/CodeReviewComments).
- Use `gofmt` for formatting your code.
- Ensure your code passes the `golint` linter checks.

## Additional Notes

- If you have any questions about how to interpret these guidelines or about contributing in general, feel free to ask in an issue or pull request.

Thank you for contributing!
