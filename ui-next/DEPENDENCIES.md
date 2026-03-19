# Pinned Dependency Versions

Verified 2026-03-19 via npm dist-tags. Agents should reference this file when adding dependencies.

> **Note:** The original spec referenced Tailwind CSS v5, Vite 6, and Storybook 8.
> These have been corrected to the actual latest stable versions below.

## Core

| Package | Version | Notes |
|---------|---------|-------|
| react | 19.2.4 | |
| react-dom | 19.2.4 | |
| typescript | 5.9.3 | |

## Build & Monorepo

| Package | Version | Notes |
|---------|---------|-------|
| vite | 8.0.1 | Spec said v6; v8 is current |
| turbo | 2.8.20 | |
| @biomejs/biome | 2.4.8 | |

## CSS & Design

| Package | Version | Notes |
|---------|---------|-------|
| tailwindcss | 4.2.2 | CSS-first `@theme` configuration |
| @tailwindcss/vite | 4.2.2 | Vite plugin for Tailwind v4 |
| shadcn | 4.1.0 | CLI — fully supports Tailwind v4 `@theme` |
| class-variance-authority | 0.7.1 | |
| clsx | 2.1.1 | |
| tailwind-merge | 3.5.0 | v3 line for Tailwind v4 compatibility |
| tw-animate-css | 1.3.4 | Replaces deprecated `tailwindcss-animate` for Tailwind v4 |

## UI Primitives

| Package | Version | Notes |
|---------|---------|-------|
| lucide-react | 0.577.0 | |
| @radix-ui/react-dialog | 1.1.15 | (and other Radix primitives as needed by shadcn) |

## State & Routing

| Package | Version | Notes |
|---------|---------|-------|
| @tanstack/react-query | 5.91.2 | |
| @tanstack/react-router | 1.167.5 | |
| zustand | 5.0.12 | |

## Charts & Editor

| Package | Version | Notes |
|---------|---------|-------|
| recharts | 3.8.0 | |
| @codemirror/view | 6.40.0 | CodeMirror 6 ecosystem |

## i18n

| Package | Version | Notes |
|---------|---------|-------|
| i18next | 25.8.20 | |
| react-i18next | 16.5.8 | |

## Testing

| Package | Version | Notes |
|---------|---------|-------|
| vitest | 4.1.0 | |
| @testing-library/react | 16.3.0 | Component testing |
| @testing-library/jest-dom | 6.6.3 | DOM assertion matchers |
| @playwright/test | 1.58.2 | |
| storybook | 10.3.1 | Spec said v8; v10 is current |

## Tooling

| Package | Version | Notes |
|---------|---------|-------|
| dependency-cruiser | 17.3.9 | |

## Protobuf

| Package | Version | Notes |
|---------|---------|-------|
| @bufbuild/protobuf | 2.11.0 | |
| @bufbuild/buf | 1.66.1 | |

## Utilities

| Package | Version | Notes |
|---------|---------|-------|
| pako | 2.1.0 | zlib decompression for API responses |
