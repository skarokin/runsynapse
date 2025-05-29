import express from "express";
import { Webhooks } from "@octokit/webhooks";
import { createClient } from "@supabase/supabase-js";
import { handleRepositoryInstallationEvent, handleInstallationEvent } from "./events/installations.ts";

// 1. go to smee.io
// 2. create new channel
// 3. copy URL and run `smee --url https://smee.io/SlIyBfMbLho7xP8T --path /webhook --port 3000`
// 4. npx tsx --env-file=.env app.ts
if (!process.env.WEBHOOK_SECRET) {
    throw new Error("WEBHOOK_SECRET is not set");
}

if (!process.env.SUPABASE_URL || !process.env.SUPABASE_SERVICE_ROLE_KEY) {
    throw new Error("Supabase URL and Anon Key must be set in environment variables");
}

const supabase = createClient(
    process.env.SUPABASE_URL,
    process.env.SUPABASE_SERVICE_ROLE_KEY    // need to bypass auth for webhooks
);

const webhooks = new Webhooks({
    secret: process.env.WEBHOOK_SECRET,
});

const app = express();

app.use('/webhook', express.raw({ type: 'application/json' }));

app.post('/webhook', async (req: express.Request, res: express.Response) => {
    try {
        const signature = req.headers["x-hub-signature-256"];
        const body = req.body.toString();
        const githubEvent = req.headers['x-github-event'] as string;

        if (!signature || typeof signature !== 'string') {
            res.status(401).send("Unauthorized");
            return;
        }

        // verify the webhook signature
        if (!(await webhooks.verify(body, signature))) {
            res.status(401).send("Unauthorized");
            return;
        }

        const payload = JSON.parse(body);

        switch (githubEvent) {
            case 'ping':
                console.log('Ping event received - webhook is working!');
                break;
            
            case 'push':
                console.log('Push event:', {
                    ref: payload.ref,
                    repository: payload.repository.full_name,
                    pusher: payload.pusher.name
                });
                // TODO: trigger Cloud Run build for this branch
                // we are given installationID in payload so we can look up the required gcp credentials in DB
                break;
            
            // note: action 'created' is handled on the sveltekit app side
            // we only need to handle 'deleted', 'suspended', 'unsuspended' and 'new_permissions_accepted' here
            case 'installation':
                console.log('Installation event');
                handleInstallationEvent(supabase, payload);
                break;
                
            case 'installation_repositories':
                console.log('Installation repositories event');
                handleRepositoryInstallationEvent(supabase, payload);
                break;
            
            default:
                console.log(`Unhandled event: ${githubEvent}`);
        }

        res.status(200).send('OK');

    } catch (error) {
        console.error('Webhook processing error:', error);
        res.status(500).send('Internal Server Error');
    }
});

// health check endpoint (don't remove req even though unused; express requires it)
app.get('/health', (req: express.Request, res: express.Response) => {
    res.status(200).json({ status: 'healthy', timestamp: new Date().toISOString() });
});

const port = process.env.PORT || 3000;

app.listen(port, () => {
    console.log(`Webhook server running on port ${port}`);
    console.log(`Webhook endpoint: http://localhost:${port}/webhook`);
    console.log(`Health check: http://localhost:${port}/health`);
});