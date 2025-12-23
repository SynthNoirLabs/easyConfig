# Task: Accessibility (A11y) Audit & Improvements

## 🎯 Objective
Ensure EasyConfig is usable by everyone, including screen reader users.

## 📝 Description
Developer tools often neglect a11y. We should be better.

## ✅ Requirements
1.  **Semantic HTML**: Replace `div`s with `button`, `nav`, `main`.
2.  **Keyboard Nav**: Ensure Tab order is logical. Focus indicators must be visible.
3.  **ARIA**: Add `aria-label` to icon-only buttons.
4.  **Monaco**: Configure Monaco Editor for accessibility (Screen Reader mode).

## 🛠️ Technical Implementation
-   **Frontend**: React component refactoring. Use `eslint-plugin-jsx-a11y`.
