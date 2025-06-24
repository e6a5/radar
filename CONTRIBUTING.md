# Contributing to Radar

Thank you for your interest in contributing to the Radar project! We welcome contributions from everyone.

## Code of Conduct

By participating in this project, you agree to abide by our Code of Conduct (see CODE_OF_CONDUCT.md).

## How to Contribute

### Reporting Bugs

Before creating bug reports, please check the existing issues to avoid duplicates. When creating a bug report, include:

- **Use a clear and descriptive title**
- **Describe the exact steps to reproduce the problem**
- **Provide specific examples** to demonstrate the steps
- **Describe the behavior you observed** and **explain what behavior you expected**
- **Include screenshots** if applicable
- **Environment details**: OS, terminal, Go version

### Suggesting Enhancements

Enhancement suggestions are welcome! Please:

- **Use a clear and descriptive title**
- **Provide a detailed description** of the suggested enhancement
- **Explain why this enhancement would be useful**
- **Include examples** of how the feature would work

### Pull Requests

1. **Fork the repository** and create your branch from `main`
2. **Make your changes** following our coding standards
3. **Add tests** for any new functionality
4. **Ensure all tests pass**: `go test ./...`
5. **Update documentation** if needed
6. **Follow the pull request template**

#### Development Setup

```bash
# Clone your fork
git clone https://github.com/YOUR_USERNAME/radar.git
cd radar

# Install dependencies
go mod download

# Run the application
go run main.go

# Run tests
go test ./...

# Build the binary
go build -o radar main.go
```

#### Coding Standards

- **Follow Go conventions**: Use `gofmt`, `golint`, and `go vet`
- **Write clear commit messages**: Use present tense ("Add feature" not "Added feature")
- **Keep functions focused**: Single responsibility principle
- **Add comments**: Especially for exported functions and complex logic
- **Handle errors appropriately**: Don't ignore errors
- **Use meaningful variable names**: Avoid abbreviations when possible

#### Testing

- Write unit tests for new functionality
- Ensure existing tests continue to pass
- Test on multiple platforms when possible
- Include integration tests for real data collection features

#### Performance Considerations

- This is a real-time application, so performance matters
- Profile your changes if they affect the main rendering loop
- Be mindful of memory allocations in hot paths
- Test with different terminal sizes

## Project Structure

```
.
├── main.go              # Application entry point
├── radar/               # Core radar package
│   ├── config.go        # Configuration management
│   ├── display.go       # Main display controller
│   ├── renderer.go      # Screen rendering logic
│   ├── signal.go        # Signal management
│   ├── realdata.go      # Real data collection
│   └── utils.go         # Utility functions
├── demo.sh              # Demo script
└── README.md            # Project documentation
```

## Feature Development Guidelines

### Signal Types
- Each signal type should have unique visual characteristics
- Movement patterns should be realistic for the signal type
- New signal types need corresponding filter controls

### Real Data Collection
- Always include fallback mechanisms for when real data isn't available
- Respect system permissions and fail gracefully
- Include timeout protection for system commands
- Test cross-platform compatibility

### UI/UX
- Maintain smooth animation (80ms refresh rate target)
- Ensure responsive design across terminal sizes
- Keep keyboard controls intuitive and documented
- Test with different color schemes and terminal types

## Release Process

1. **Update version** in relevant files
2. **Update CHANGELOG.md** with new features and fixes
3. **Create a release tag**: `git tag v1.x.x`
4. **Push tag**: `git push origin v1.x.x`
5. **GitHub Actions** will handle the rest

## Getting Help

- **Documentation**: Check the README.md first
- **Issues**: Search existing issues or create a new one
- **Discussions**: Use GitHub Discussions for questions
- **Real-time help**: Check if there are any community chat channels

## Recognition

Contributors will be recognized in:
- The project README
- Release notes for significant contributions
- Special thanks for major features or bug fixes

## License

By contributing to this project, you agree that your contributions will be licensed under the MIT License.

---

Thank you for making Radar better! 🎯 