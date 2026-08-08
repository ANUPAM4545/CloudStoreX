export default function AICenter() {
  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-6">AI Storage Intelligence</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold mb-2">AI Assistant</h3>
          <p className="text-sm text-gray-500 mb-4">Chat with the intelligent operations assistant.</p>
          <a href="/dashboard/ai/assistant" className="text-blue-600 hover:underline">Open Assistant</a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold mb-2">Cost Optimizer</h3>
          <p className="text-sm text-gray-500 mb-4">View AI-generated cost reduction recommendations.</p>
          <a href="/dashboard/ai/cost" className="text-blue-600 hover:underline">View Insights</a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold mb-2">Capacity Forecast</h3>
          <p className="text-sm text-gray-500 mb-4">Predictive models for storage growth and reliability.</p>
          <a href="/dashboard/ai/forecasts" className="text-blue-600 hover:underline">View Forecasts</a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold mb-2">Security Intelligence</h3>
          <p className="text-sm text-gray-500 mb-4">Anomaly detection and AI-scored risk analysis.</p>
          <a href="/dashboard/ai/security" className="text-blue-600 hover:underline">View Intelligence</a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold mb-2">Policy Advisor</h3>
          <p className="text-sm text-gray-500 mb-4">Automated routing policy optimizations.</p>
          <a href="/dashboard/ai/policies" className="text-blue-600 hover:underline">View Recommendations</a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold mb-2">AI Settings</h3>
          <p className="text-sm text-gray-500 mb-4">Manage Provider Abstractions, Context Memory, and Cache.</p>
          <a href="/dashboard/ai/settings" className="text-blue-600 hover:underline">Configure</a>
        </div>
      </div>
    </div>
  );
}
