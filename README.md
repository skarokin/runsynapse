# runsynapse

Capture half-formed thoughts on the run, let AI connect the dots. Inspired by personal Discord server workflows.

## Tech Stack
- **Go + API Gateway + Lambda** - Serverless backend
- **SQS + Lambda + Gemini** - Async vector embedding processing
- **PostgreSQL + pgvector + Gemini** - Hybrid search, thought connection, and summaries
- **SvelteKit + Cloudflare** - Fast frontend hosting
- **S3** - Object storage