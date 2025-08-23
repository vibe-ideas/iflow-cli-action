# TOC Generation Example

This example demonstrates how to automatically generate a table of contents for all Markdown files in your repository.

## How It Works

The workflow uses the `markdown-toc` npm package to automatically generate and insert a table of contents into all Markdown files. It runs on:

1. Pushes to the main branch that modify Markdown files
2. Pull requests to the main branch that modify Markdown files
3. Manual triggering via workflow_dispatch

The workflow has two main steps:
1. Add TOC markers (`<!-- toc -->` and `<!-- tocstop -->`) to all Markdown files that don't already have them
2. Generate and insert the actual table of contents between these markers

## Setup Instructions

1. Add this workflow file to your repository at `.github/workflows/generate-toc.yml`
2. Ensure your repository has the proper permissions for GitHub Actions to commit changes
3. The workflow will automatically process all Markdown files in your repository

## Configuration Options

You can customize the workflow by modifying these parameters:

- `branches`: Specify which branches should trigger the workflow
- `paths`: Filter which file changes should trigger the workflow
- `node-version`: Specify which Node.js version to use
- `commit_message`: Customize the commit message for TOC updates

## Example Output

After running, your Markdown files will have a table of contents automatically inserted. For example:

```markdown

<!-- toc -->

- [Section 1](#section-1)
  * [Subsection 1.1](#subsection-11)
- [Section 2](#section-2)

<!-- tocstop -->

# Section 1

## Subsection 1.1

Content here...

# Section 2

More content...
```

The TOC is automatically updated whenever Markdown files are modified and pushed to the repository.