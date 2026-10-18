# FlowForge hands-on

This repository is a small shop platform with three services in three languages:

| Path | Stack | What it is |
|---|---|---|
| `services/orders-go` | Go 1.23, Dockerfile | Orders API: tests, `go vet`, gofmt, container image |
| `services/pricing-py` | Python 3.12, uv | Discount rules: ruff and pytest |
| `web/storefront` | Node 22, npm | Cart logic: `node --test`, build, and a demo deploy script |

None of the files contain FlowForge configuration. FlowForge works out the pipelines from the code itself.

You will:
1. Push this repository to GitHub.
2. Let FlowForge plan its CI in the portal.
3. Publish a draft PR from the portal and watch the real runs go green.
4. Change a step from the UI and republish.
5. Merge, then deploy to staging and production with environments and an approval gate. This needs no
   cloud account.

Time: about 30 minutes. Every command runs in your terminal on this Mac.

---

## Step 0: Check your tools (1 min)

```sh
gh auth status            # should say: Logged in to github.com account devopsmaker
git --version
```

If `gh` isn't logged in, run `gh auth login` first. FlowForge uses this login to read your repository and
open the PR.

## Step 1: Put this code on GitHub (2 min)

```sh
cd /Users/opstree/Documents/github.io/flowforge-demo
git init -b main
git add .
git commit -m "Shop platform: orders (Go), pricing (Python), storefront (Node)"
gh repo create flowforge-demo --public --source . --push
```

**Why public?** GitHub only allows required reviewers on environments in private repos on paid plans.
Step 7 uses a reviewer to gate production. Nothing in this repo is secret. If you prefer `--private`,
everything works except that approval gate.

Check it: `gh repo view --web` opens the repository. Its **Actions** tab is empty, because there are no
workflows yet.

## Step 2: Start the portal (1 min)

```sh
cd /Users/opstree/Documents/github.io/FlowForge
uv run flowforge ui
```

Your browser opens <http://127.0.0.1:8765>. In the sidebar, the dots for **OPA**, **actionlint** and
**zizmor** should be green. Those are the policy engine and the two workflow validators every plan
passes through.

Leave this terminal running and use a second terminal for later commands.

## Step 3: Plan the pipeline (3 min)

1. On the **Plan** page, pick **GitHub repo** and enter `devopsmaker/flowforge-demo`.
2. Click **Analyze**. It takes a few seconds: FlowForge clones the repo into a temporary directory and
   reads it.

You should see:
- **3 components**: `services-orders-go`, `services-pricing-py`, `web-storefront`
- **Policy: pass**, with one warning (REPO001, which you'll fix in step 8)
- **Validation: 3/3** (YAML, actionlint, zizmor)
- **4 questions** in the right-hand panel

Explore the tabs before answering:

| Tab | Look for |
|---|---|
| **Pipeline** | Five CI jobs: a dependency and secret scan, one job per service, and a container build for the Go service. Click a step to see why it's there |
| **Workflows** | `flowforge-ci.yml`. Every `# why [L-…]` comment is a link to the decision behind that line, for example why Go 1.23 (`go.mod:3`) or why Node 22 (`.nvmrc:1`) |
| **Decision Ledger** | Every value with its source (detected, default, policy, you), confidence and the file and line it came from |
| **Evidence** | The raw facts the detectors found |

## Step 4: Answer the questions (2 min)

In the **Questions** panel:

| Question | Answer | Why |
|---|---|---|
| `services-orders-go.deploy.target` | `none` | CI only for now |
| `services-pricing-py.deploy.target` | `none` | CI only for now |
| `web-storefront.deploy.target` | `none` for now | You'll add a deploy in step 7 |
| `services-orders-go.ci.step.format_check` | click **Accept best guess** (`gofmt`) | FlowForge asks because a format check fails CI on unformatted code. This code is formatted |

Click **Apply 4 answers & regenerate**. The plan regenerates in under a second and highlights what changed. The
status should now read **Ready**.

## Step 5: Go-live check (1 min)

Open the **Go-live** tab. It reads your repository's GitHub settings live: whether Actions is enabled and
every pinned action is allowed. All checks should show **ready**. If one doesn't, its **Fix on GitHub**
button changes exactly that setting.

## Step 6: Publish and watch real CI (5 min)

1. Click **Open draft PR**, read the dialog, then click **Write to GitHub**.
2. A **Live CI** card appears under the header. It follows the run on GitHub's own runners:
   - *Waiting for GitHub to start FlowForge CI…*
   - *FlowForge CI is running…*, with each job and its current step
   - *FlowForge CI passed on GitHub's runners* (about 2 to 4 minutes; the dependency scan downloads its
     vulnerability database first)
3. Click the PR number to see the same run on GitHub: the PR adds `.github/workflows/flowforge-ci.yml`,
   `.flowforge/config.yaml` (your answers) and `.flowforge/pim.lock.json` (every decision and the hash of
   each workflow).

**If a job fails,** the card shows the failing step, whose problem it is (pipeline, repository settings,
or your project's code) and the end of the log. Note what it says. That diagnosis is part of the tool.

## Step 7: Change a step from the UI, then republish (3 min)

1. Open the **Decision Ledger** tab and find `web-storefront.ci.step.unit_test`. Click it.
2. In **Override this value**, enter `npm test -- --test-reporter=spec` and apply.
3. The plan regenerates. In **Workflows**, the changed line is highlighted, and its `why` now says the value
   came from you.
4. Click **Open draft PR** again. The same PR is updated with a new commit, and Live CI follows the new run.

Free-form commands are accepted only because a human typed them. Detection and the LLM can never
introduce one, and policy enforces that.

## Step 8: Merge, and add yourself as workflow owner (2 min)

In the second terminal:

```sh
cd /Users/opstree/Documents/github.io/flowforge-demo
gh pr ready flowforge/onboard       # mark the draft ready
gh pr merge flowforge/onboard --squash --delete-branch
git pull
mkdir -p .github && echo "/.github/ @devopsmaker" > .github/CODEOWNERS
git add .github/CODEOWNERS && git commit -m "Own the workflows" && git push
```

The push to `main` runs FlowForge CI again: **Actions** tab → *FlowForge CI*. Analyze again in the portal and
the REPO001 warning is gone.

Prove the workflows are exactly what the reviewed lockfile renders. This needs no network and no
FlowForge service:

```sh
cd /Users/opstree/Documents/github.io/FlowForge
uv run flowforge verify /Users/opstree/Documents/github.io/flowforge-demo    # verify: OK
```

## Step 9: Deploy to staging and production, with no cloud account (8 min)

`web/storefront/scripts/deploy.sh` is a stand-in deploy. It builds the storefront and prints where it
"deployed". It refuses to run without a `DEMO_TOKEN` secret, so you can see secrets wired per environment.

1. In the portal, click **Analyze** again (the repo now has the merged workflows), then answer:

   | Question | Answer |
   |---|---|
   | `web-storefront.deploy.target` | `command` |
   | `web-storefront.deploy.environments` | `staging, production` (best guess) |
   | `web-storefront.deploy.deploy_command` | `./scripts/deploy.sh` |
   | `web-storefront.deploy.secret_names` | `DEMO_TOKEN` |
   | the rest (rollback, layout, GitHub Environments) | **Accept best guess** on each |

   A new `flowforge-deploy-web-storefront.yml` appears. Staging deploys on every push to `main`, and
   production on a published release. Policy warns with PIM011 that this target deploys from source: an
   honest note, not a failure.

2. Open **Go-live**. Now it has work for you:
   - **Environment `staging`** and **Environment `production`**: missing. Click **Fix on GitHub** on each.
     Production is created with you as required reviewer, limited to `v*` tags.
   - **Secret `DEMO_TOKEN`** for each environment: type any value (for example `hello-staging`) and click
     **Set secret**. It is sealed by the GitHub CLI; FlowForge never stores it.
   - Click **Re-check**: everything is **ready**, except one item listed as *outside GitHub*.

3. Click **Open draft PR** and wait for Live CI to go green. Then merge it as in step 8 (`gh pr ready …`,
   `gh pr merge …`).

4. The merge is a push to `main`, so **deploy to staging** runs. Watch it in **Actions**: the log ends with
   `Deployed storefront <sha> to staging (token length 13)`.

5. Promote to production:

   ```sh
   cd /Users/opstree/Documents/github.io/flowforge-demo
   git pull && gh release create v0.1.0 --title v0.1.0 --notes "First release"
   ```

   The production deploy **waits for your approval**: **Actions** → the run → *Review deployments* →
   approve. Then it deploys the released commit.

## Step 10: Break it on purpose (optional, 3 min)

See how failures are reported, locally before any push:

```sh
cd /Users/opstree/Documents/github.io/flowforge-demo
sed -i '' 's/got != 750/got != 751/' services/orders-go/orders_test.go
cd ../FlowForge && uv run flowforge run-ci ../flowforge-demo --component services-orders-go
```

The output shows `[FAIL] Test` with the Go test failure, and `Build` not run. `flowforge run-ci` runs
the generated CI commands on a throwaway copy of the repo, on your machine. Undo the change with
`git -C ../flowforge-demo checkout services/orders-go/orders_test.go`.

To see the same failure on GitHub, push the change on a branch, open a PR, and read the failing job's log
under the PR's **Checks** tab.

## Clean up

```sh
gh repo delete devopsmaker/flowforge-demo --yes   # needs: gh auth refresh -s delete_repo
```

## What to tell me afterwards

Anything that didn't match this guide: the step, what you expected, what you saw, and a screenshot of the
Live CI card or the failing log. Those are the missing pieces of real-runner evidence for FlowForge.
