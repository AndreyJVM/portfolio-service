---
name: english-code-comments
description: Standardized rule and playbook for writing, editing, and maintaining all code comments, docstrings, and inline annotations strictly and exclusively in English across all codebases in the repository.
version: 1.0.0
tags: [golang, coding-standards, comments, documentation, agentic-skill]
---

# Skill Playbook: English Code Comments

This skill defines the standard requirement for all source code comments across the repository.

## 1. Rule: English Only for Code Comments

All comments, docstrings, and annotations within source code files (Golang, Shell, Python, JavaScript, CSS/HTML, Dockerfile, etc.) must be written **exclusively in English**.

### Scope

1. **Doc Comments & Package Comments:**
   - Package-level comments (`// Package ...`).
   - Exported and unexported functions, methods, structs, interfaces, and constants.
   - Follow standard language conventions (e.g., Go doc format starting with the identifier name: `// GetLesson loads...`).

2. **Inline Comments:**
   - Explanations of complex logic, algorithms, concurrency patterns, or non-obvious design choices.
   - Rationale for edge-case handling or defensive programming measures (e.g., DoS protection, bounds checking, timeouts).

3. **Test Comments:**
   - Descriptions of test steps, assertions, mock setups, and table-driven test cases (`// 1. Setup mock server`, `// 2. Execute request`).

4. **TODOs, FIXMEs, and Annotations:**
   - Tags like `// TODO:`, `// FIXME:`, `// NOTE:`, `// BUG:`.

## 2. Distinction Between Comments and Data

- **Code comments (`// ...`, `/* ... */`, `# ...`)**: Must always be in English.
- **Localized user-facing content**: String literals, HTML/template content, translation dictionaries, or course materials intended for users may remain in Russian or the designated target locale. Do not confuse code comments with user-facing string literals.

## 3. Best Practices

- Use clear, natural, and idiomatic technical English.
- Keep comments concise and focused on *why* rather than *what* when the code itself is self-explanatory.
- When refactoring or updating existing files with legacy non-English comments, translate them to English as part of the modification.
