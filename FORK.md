# Fork of linden-project/linny-mcp-server

This is [torreirow](https://github.com/torreirow)'s fork of
[linden-project/linny-mcp-server](https://github.com/linden-project/linny-mcp-server).

Forked from: v0.2.0 (ee628a0)

## Workflow

- `main` is the fork's own integration branch. It is what torreirow's
  setup pins; it is never proposed upstream as a whole.
- The `upstream` remote points at `linden-project/linny-mcp-server`.
- Every change lives on its own topic branch off `upstream/main` and is merged
  into `main`. Each change keeps the current upstream behaviour as the default
  (new behaviour is opt-in through configuration), so the branch can later be
  offered upstream as a PR unchanged.
- No upstream PRs for now; whether to offer a topic branch is decided per branch.

Sync with upstream:

```sh
git fetch upstream
git merge upstream/main   # on main
```

## Deviations from upstream

| Topic branch | Merged in main | Upstream PR | Description |
|--------------|----------------|-------------|-------------|
| –            | –              | –           | None yet    |
