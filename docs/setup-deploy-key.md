Absolutely! Here’s a **minimal, step-by-step markdown doc** to help you remember how to set up and use a GitHub Deploy Key for Terragrunt private module access between two repos:

---

# Setting Up a GitHub Deploy Key for Terragrunt Module Access

This guide walks you through giving a “live” repo CI/CD access to private modules stored in a “stack” repo using a deploy key.

---

## **Steps**

### 1. **In the Stack Repo (module source)**

Make sure you are in the directory of your *stack* repo (e.g. `terragrunt-template-stack`).

```sh
# Get the repo name in owner/name format (requires gh CLI)
REPO_STACK=$(gh repo view --json nameWithOwner --jq .nameWithOwner)
echo $REPO_STACK
```

Create an environment variable for the name of the key:
```sh
KEYNAME="DEPLOY_KEY_TG_STACK"
KEYPATH="/tmp/$KEYNAME"
```

Generate an SSH keypair **(no passphrase)**:

```sh
ssh-keygen -t ed25519 -C "$KEYNAME" -f "$KEYPATH" -N ""
```

Add the **public key** as a deploy key to the *stack* repo (read-only):

```sh
gh api \
  --method POST \
  -H "Accept: application/vnd.github+json" \
  "/repos/$REPO_STACK/keys" \
  -f title="$KEYNAME" \
  -f key="$(cat $KEYPATH.pub)" \
  -F read_only=true
```

Open the deploy key url on your browser to ensure it has been created:
```sh
echo "https://github.com/$REPO_STACK/settings/keys"
```

### 2. **In the Live Repo (consumer / workflow)**

Get the full `owner/name` of your *live* repo (where you’ll run Terragrunt):

```sh
REPO_LIVE=$(gh repo view --json nameWithOwner --jq .nameWithOwner)
echo $REPO_LIVE
```

Add the **private key** (`$KEYPATH`) as a repository secret (example: `DEPLOY_KEY_TG_STACK`) to the *live* repo:

```sh
gh secret set DEPLOY_KEY_TG_STACK --repo "$REPO_LIVE" < "$KEYPATH"
```

Ensure the secret exists in the live repository under `Repository secrets` by cliking on this link:
```sh
echo "https://github.com/$REPO_LIVE/settings/secrets/actions"
```

### 3. **Use the Deploy Key in Your Workflow**

In your workflow YAML in the *live* repo, add before any Terragrunt step:

```yaml
- name: Set up SSH for Terragrunt module pulls
  uses: webfactory/ssh-agent@v0.9.0
  with:
    ssh-private-key: ${{ secrets.DEPLOY_KEY_TG_STACK }}
```

Your Terragrunt `source` should use the SSH format:

```hcl
source = "git::git@github.com:ConsciousML/terragrunt-template-stack.git//units/foo?ref=..."
```

---

## **Notes**

* **Do NOT use the same deploy key in more than one GitHub repo.** Generate a new key for each live-stack pairing.
* **Do NOT commit private keys to version control.**
* **You only need to do this once per repo pairing, unless rotating keys.**

---

*You now have automated, secure, repo-scoped access for Terragrunt to private modules in CI!*
