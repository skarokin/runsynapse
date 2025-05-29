<script>
    import { Button } from "$lib/components/ui/button/index.js";
</script>

<div class="space-y-2 w-full sm:w-3/4">
    <h1 class="text-4xl font-bold leading-tight">Privacy Policy</h1>
    <h3 class="text-xl text-muted-foreground text-pretty">
        Last updated: 05/25/2025
    </h3>
</div>

<div class="grid grid-cols-2 gap-x-16 gap-y-6 w-full sm:w-3/4">
    <div class="space-y-2">
        <h3 class="font-medium">Information We Collect</h3>
        <p class="text-sm text-muted-foreground">
            When you use runsynapse, we collect your Google account information (email, profile, Cloud project access), GitHub integration data (repository names, commits, webhooks), and usage data (deployment logs, performance metrics).
        </p>
    </div>

    <div class="space-y-2">
        <h3 class="font-medium">How We Use Your Information</h3>
        <p class="text-sm text-muted-foreground">
            We use your data to deploy and manage applications on Google Cloud Run, monitor deployment status, process GitHub webhooks for automated deployments, and provide customer support.
        </p>
    </div>

    <div class="space-y-2">
        <h3 class="font-medium">Data Security</h3>
        <p class="text-sm text-muted-foreground">
            All data is encrypted in transit and at rest. Access tokens are securely stored using Supabase. We perform regular security audits and never log sensitive credentials.
        </p>
    </div>

    <div class="space-y-2">
        <h3 class="font-medium">Data Sharing</h3>
        <p class="text-sm text-muted-foreground">
            We don't sell or share your personal information with third parties, except with Google Cloud and GitHub APIs to provide our services, or when required by law.
        </p>
    </div>

    <div class="space-y-2">
        <h3 class="font-medium">Your Rights</h3>
        <p class="text-sm text-muted-foreground">
            You can access, delete, or export your data at any time. You can disconnect Google Cloud or GitHub integrations from your dashboard and request account deletion.
        </p>
    </div>

    <div class="space-y-2">
        <h3 class="font-medium">Data Retention</h3>
        <p class="text-sm text-muted-foreground">
            We retain your data while your account is active. Upon deletion, personal information is removed within 30 days, deployment logs within 90 days.
        </p>
    </div>

    <div class="space-y-2">
        <h3 class="font-medium">Third-Party Services</h3>
        <p class="text-sm text-muted-foreground">
            We integrate with
            <a
                href="https://policies.google.com/privacy"
                class="text-primary hover:underline"
                target="_blank"
                rel="noopener noreferrer"
            >
                Google Cloud
            </a>, 
            <a
                href="https://docs.github.com/en/site-policy/privacy-policies/github-privacy-statement"
                class="text-primary hover:underline"
                target="_blank"
                rel="noopener noreferrer"
            >
                GitHub
            </a>, and
            <a
                href="https://supabase.com/privacy"
                class="text-primary hover:underline"
                target="_blank"
                rel="noopener noreferrer"
            >
                Supabase
            </a>. Please review their privacy policies.
        </p>
    </div>

    <div class="space-y-2">
        <h3 class="font-medium">Contact Us</h3>
        <p class="text-sm text-muted-foreground">
            Questions about this Privacy Policy? Contact <strong>sean.kim0327@gmail.com</strong>
        </p>
    </div>
</div>

<Button href="/">Back to Home</Button>
