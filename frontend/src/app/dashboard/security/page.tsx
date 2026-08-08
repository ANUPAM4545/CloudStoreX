export default function SecurityCenter() {
  return (
    <div className="p-8">
      <h1 className="text-3xl font-bold mb-6">Security Center</h1>
      
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold mb-2">Organizations & Teams</h3>
          <p className="text-sm text-gray-500 mb-4">Manage enterprise structure, departments, and user groups.</p>
          <a href="/dashboard/security/organizations" className="text-blue-600 hover:underline">Manage</a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold mb-2">Access Policies</h3>
          <p className="text-sm text-gray-500 mb-4">Configure RBAC roles, ABAC rules, and Zero Trust settings.</p>
          <a href="/dashboard/security/access-policies" className="text-blue-600 hover:underline">Configure</a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold mb-2">Active Sessions</h3>
          <p className="text-sm text-gray-500 mb-4">Monitor and revoke active user and device sessions.</p>
          <a href="/dashboard/security/sessions" className="text-blue-600 hover:underline">View Sessions</a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold mb-2">API Keys</h3>
          <p className="text-sm text-gray-500 mb-4">Manage programmatic access tokens and scopes.</p>
          <a href="/dashboard/security/api-keys" className="text-blue-600 hover:underline">Manage Keys</a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold mb-2">Service Accounts</h3>
          <p className="text-sm text-gray-500 mb-4">Manage machine identities and secret rotation.</p>
          <a href="/dashboard/security/service-accounts" className="text-blue-600 hover:underline">Manage Accounts</a>
        </div>

        <div className="bg-white p-6 rounded-lg shadow-sm border border-gray-100">
          <h3 className="text-lg font-semibold mb-2">Audit Logs</h3>
          <p className="text-sm text-gray-500 mb-4">View immutable security events and alerts.</p>
          <a href="/dashboard/security/audit" className="text-blue-600 hover:underline">View Logs</a>
        </div>
      </div>
    </div>
  );
}
