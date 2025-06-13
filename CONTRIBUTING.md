# Contributing to Vibecast

Thank you for your interest in contributing to Vibecast!  
We welcome all kinds of contributions: bug reports, feature requests, code, documentation, and design.

## How to Contribute

1. **Open an Issue**  
   - If you find a bug or have a feature request, please [open an issue](https://github.com/pedrobarco/vibecast/issues).
   - Search existing issues before creating a new one.

2. **Fork the Repository**  
   - Click "Fork" at the top right of the [repo page](https://github.com/pedrobarco/vibecast).

3. **Clone Your Fork**

   ```bash
   git clone https://github.com/yourusername/vibecast.git
   cd vibecast
   ```

4. **Create a Branch**

   ```bash
   git checkout -b my-feature
   ```

5. **Make Your Changes**
   - Follow Go best practices and keep code idiomatic.
   - Write clear commit messages.
   - Add or update tests if appropriate.

6. **Run Tests**

   ```bash
   go test ./...
   ```

7. **Push and Open a Pull Request**

   ```bash
   git push origin my-feature
   ```

   - Go to your fork on GitHub and open a pull request (PR) to `main`.

## Code Style

- Use `gofmt` and `goimports` before submitting.
- Keep functions small and focused.
- Document exported functions and types.
- Use descriptive variable and function names.

## Commit Messages

- Use [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/):
  - `feat: ...` for new features
  - `fix: ...` for bug fixes
  - `docs: ...` for documentation
  - `refactor: ...` for code refactoring
  - `test: ...` for tests
  - `chore: ...` for maintenance

## Code of Conduct

Be respectful and inclusive. See [Contributor Covenant](https://www.contributor-covenant.org/) for guidance.

## Questions?

Open an issue or start a discussion!

Thank you for helping make Vibecast better!
