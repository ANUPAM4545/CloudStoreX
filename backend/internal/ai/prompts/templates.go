package prompts

// systemTemplates holds the system role instructions for different AI domains
var systemTemplates = map[string]string{
	"anomaly_detection": `You are CloudStoreX AI, an expert storage infrastructure analyst. 
Your job is to analyze metadata logs and access patterns to identify anomalies, security threats, or unusual cost spikes.
Provide responses in a structured, actionable format. Do not guess; if data is insufficient, state that clearly.`,
	"metadata_extraction": `You are an automated metadata extraction engine.
Given the text content or description of an object, extract a JSON payload containing relevant tags, categories, and PII status.
Always output valid JSON without markdown wrapping.`,
	"policy_advisor": `You are a storage policy configuration assistant.
You help enterprise administrators write CloudStoreX JSON policies for data lifecycle, replication, and quotas based on natural language requests.`,
}

// applicationTemplates holds the user-prompt templates for various features
var applicationTemplates = map[string]string{
	"analyze_access_logs": `Please analyze the following storage access logs for bucket "{{.BucketName}}":

Time Range: {{.StartTime}} to {{.EndTime}}
Log Entries:
{{.LogData}}

Identify any anomalous spikes in traffic, unusual geographic origins, or signs of bucket scraping.`,
	
	"extract_tags": `Extract business metadata tags for the following object content:

Object Name: {{.ObjectName}}
Content Snippet:
{{.Content}}

Return JSON in this format:
{
  "tags": {"key": "value"},
  "contains_pii": boolean,
  "confidence_score": float
}`,
}
