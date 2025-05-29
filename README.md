# runsynapse
Streamline your containerized Cloud Run workloads with better monitoring, seamless integrations, and automated workflows.
- Better UI for metrics & billing visualization.
- Easy alert setup and monitoring for proactive issue resolution.
- Direct Cloud Run integration with seamless container view switching.
- GitHub CI/CD integration with preview deployments and secrets management.

Tasks:
- [x] Google Sign-in from Supabase
- [x] GitHub App installation & connection of installation to userID
- [x] Push repository events as webhooks
- [x] Listen for updated GitHub App repository access & add/update/delete from DB accordingly
- [x] Connect GCP account and list projects
- [x] Encrypt refresh token at rest (already handled by Supabase I just learned :D)
- [x] Listen for push events from GitHub
- [ ] Figure out a good UI and user flow for adding repos, connecting to projects, viewing deployments etc! Currently horrible omg
- [ ] Dockerfile path configuration
- [ ] Link repository to GCP project
- [ ] Provision GCP resources for the user
- [ ] Set up build pipeline (either Cloud Build or clone + Dockerize + push from my VM)
- [ ] CI/CD after pushes
- [ ] Rollback, shut down, suspend, and manual deployments
- [ ] Preview deployments that shut down/spin up on branch close/creation
- [ ] Logs & metrics streaming from Cloud Run services
- [ ] Billing information
- [ ] Secrets management
- [ ] GitHub App pushes deployment status events & gives link
- [ ] Set up alerting via API (no need to handle emails and stuff myself)