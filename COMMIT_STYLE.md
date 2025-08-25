# Commit Message Style Guide

## Format

All commit messages should follow the conventional commits format:

```
<type>: <description>

[optional body]

[optional footer(s)]
```

## Types

- **feat**: A new feature
- **fix**: A bug fix
- **docs**: Documentation only changes
- **style**: Changes that do not affect the meaning of the code (white-space, formatting, etc)
- **refactor**: A code change that neither fixes a bug nor adds a feature
- **perf**: A code change that improves performance
- **test**: Adding missing tests or correcting existing tests
- **build**: Changes that affect the build system or external dependencies
- **ci**: Changes to CI configuration files and scripts
- **chore**: Other changes that don't modify src or test files

## Rules

1. **Language**: All commit messages must be in English
2. **Tense**: Use present tense ("add feature" not "added feature")
3. **Case**: Start with lowercase (except for proper nouns)
4. **Length**: Keep the description under 72 characters
5. **Description**: Be clear and concise about what changed and why

## Examples

### Good Examples
- `feat: add support for custom SQL query validation`
- `fix: correct false positive in SELECT * detection`
- `docs: update README with installation instructions`
- `ci: split cache configuration to prevent conflicts`
- `test: add test cases for multiline SQL strings`
- `refactor: extract common constants to reduce duplication`

### Bad Examples
- `Fixed bug` (missing type prefix, past tense)
- `feat: Добавлена поддержка...` (not in English)
- `update` (too vague, missing type)
- `FEAT: ADD NEW FEATURE` (all caps)
- `feat: this commit adds a new feature that allows users to...` (too long)

## Commit History Consistency

The existing commit history follows these patterns:
- Early commits may use simple descriptions (`Add test`, `Fix tests`)
- Recent commits use the conventional format (`feat:`, `fix:`, `docs:`)
- All commits should be in English
- Focus on what changed, not how it was changed
