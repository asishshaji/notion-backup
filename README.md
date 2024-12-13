# notion-backup
Notion backup using go

##` Description

This is a simple tool to keep backups of your notion workspace to a github repo. Inspired from [Notion Backup](https://github.com/darobin/notion-backup). It is designed to work as part of a github workflow.

More information about how to get the tokens are in the above link.

```yml
name: "Notion backup"

on:
  push:
    branches:
      - main
  schedule:
    - cron: "30 23 * * *"

  # Allows you to run this workflow manually from the Actions tab
  workflow_dispatch:

jobs:
  backup:
    runs-on: ubuntu-latest
    name: Backup
    timeout-minutes: 15
    steps:
      - uses: actions/checkout@v3

      - name: Run backup
        run: ./notion-backup -html -markdown
        env:
          NOTION_TOKEN: ${{ secrets.NOTION_TOKEN }}
          NOTION_FILE_TOKEN: ${{ secrets.NOTION_FILE_TOKEN }}
          NOTION_SPACE_ID: ${{ secrets.NOTION_SPACE_ID }}

      - name: Commit changes
        run: |
          git config user.name github-actions
          git config user.email github-actions@github.com
          git add .
          git commit -m "Automated snapshot"
          git push
```
