import { Logging, type protos } from '@google-cloud/logging';
import { ServicesClient } from '@google-cloud/run';
import { ServiceUsageClient } from '@google-cloud/service-usage';

export async function fetchCloudRunLogs(
    credentials: any,
    projectId: string,
    serviceName: string,
    branchName: string = "main",
    page: string | null = null,
    maxAgeDays?: string
) {
    if (!credentials || !projectId || !serviceName) {
        throw new Error('Missing required parameters for fetching Cloud Run logs');
    }

    const pageSize = 10;

    const timeFilter = maxAgeDays || new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString();

    try {
        const logging = new Logging({
            credentials: credentials,
            projectId: projectId
        });

        const filter = `resource.type="cloud_run_revision" AND resource.labels.service_name="${serviceName}-${branchName}" AND timestamp >= "${timeFilter}"`;

        const requestOptions: any = {
            resourceNames: [`projects/${projectId}`],
            filter: filter,
            orderBy: 'timestamp desc',
            pageSize: pageSize,
        };

        // only add pageToken if it exists and is not empty
        if (page && page.trim() !== '') {
            requestOptions.pageToken = page;
        }

        const response = await logging.getEntries(requestOptions);

        const logEntries = response[0]; // entries array
        const nextPageToken = response[1]?.pageToken; // next page token from the request object

        const logs = logEntries.map((entry: any) => {
            let processedData = entry.data;
            const payloadType = entry.metadata.payload;

            if (payloadType === 'textPayload') {
                // regular text log
                processedData = entry.data;
            } else if (payloadType === 'payloadNotSet' || !payloadType) {
                // no specific payload type set, return the http request (if available)
                if (entry.metadata.httpRequest) {
                    const req = entry.metadata.httpRequest;
                    processedData = `${req.requestMethod} ${req.status} ${req.responseSize || 0} Bytes ${req.latency} ${req.requestUrl} ${req.userAgent || ''}`;
                } else {
                    processedData = 'No payload data';
                }
            } else if (payloadType === 'protoPayload') {
                // proto payload, must be decoded
                processedData = `[Proto] ${entry.data?.type_url || 'Binary data'}`;
            }

            return {
                timestamp: entry.metadata.timestamp,
                severity: entry.metadata.severity,
                payloadType: payloadType,
                data: processedData
            };
        });

        return {
            logs,
            nextPageToken: nextPageToken || null,
            timeFilter: timeFilter,
        }
    } catch (error) {
        console.error('Error fetching logs:', error);
        throw error;
    }
}

export async function fetchCloudRunMetrics(credentials: any, projectId: string, serviceName: string) {
    try {
        const client = new ServicesClient({
            credentials: credentials
        });

        // first, try to list all services to find the correct one
        const parent = `projects/${projectId}/locations/-`;
        const [services] = await client.listServices({
            parent: parent
        });

        const service = services.find((s: any) =>
            s.name?.includes(serviceName) ||
            s.name === serviceName
        );

        if (!service) {
            console.log('Available services:', services.map((s: any) => s.name));
            return null;
        }

        const shortName = service.name?.split('/').pop();

        return {
            name: shortName,
            uri: service.uri,
            url: service.urls?.[0] || service.uri, // fallback to first URL if available
            ready: service.terminalCondition?.state === 'CONDITION_SUCCEEDED',
            traffic: service.traffic || [],
            lastModifier: service.lastModifier,
            generation: service.generation
        };

    } catch (error) {
        console.error('Error fetching metrics:', error);
        throw error;
    }
}

export async function enableCloudRunAPI(credentials: any, projectId: string) {
    try {
        const client = new ServiceUsageClient({
            credentials: credentials
        });

        const request = {
            name: `projects/${projectId}/services/run.googleapis.com`
        };

        const [operation] = await client.enableService(request);
        console.log('Enable API operation:', operation);

        return operation;
    } catch (error) {
        console.error('Error enabling Cloud Run API:', error);
        throw error;
    }
}