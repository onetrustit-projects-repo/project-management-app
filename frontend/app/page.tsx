import Link from 'next/link'

export default function Home() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-slate-900 to-slate-800">
      {/* Hero Section */}
      <header className="px-6 py-8">
        <nav className="max-w-6xl mx-auto flex justify-between items-center">
          <h1 className="text-2xl font-bold text-white">ProjectFlow</h1>
          <div className="flex gap-4">
            <Link href="/dashboard" className="px-4 py-2 text-white hover:text-slate-300">Dashboard</Link>
            <Link href="/projects" className="px-4 py-2 text-white hover:text-slate-300">Projects</Link>
            <Link href="/tasks" className="px-4 py-2 text-white hover:text-slate-300">Tasks</Link>
          </div>
        </nav>
      </header>

      <main className="max-w-6xl mx-auto px-6 py-20">
        <div className="text-center mb-16">
          <h2 className="text-5xl font-bold text-white mb-6">
            AI-Powered Project Management
          </h2>
          <p className="text-xl text-slate-300 mb-8 max-w-2xl mx-auto">
            Streamline your workflow with intelligent automation, real-time collaboration, 
            and powerful analytics. Built for teams who move fast.
          </p>
          <div className="flex gap-4 justify-center">
            <Link href="/dashboard" className="px-8 py-3 bg-indigo-600 text-white rounded-lg font-medium hover:bg-indigo-700 transition">
              Get Started
            </Link>
            <Link href="/projects" className="px-8 py-3 border border-slate-500 text-white rounded-lg font-medium hover:bg-slate-800 transition">
              View Demo
            </Link>
          </div>
        </div>

        {/* Features Grid */}
        <div className="grid md:grid-cols-3 gap-8 mb-16">
          <div className="bg-slate-800/50 rounded-xl p-6 border border-slate-700">
            <div className="text-3xl mb-4">📋</div>
            <h3 className="text-xl font-semibold text-white mb-2">Kanban Boards</h3>
            <p className="text-slate-400">Drag-and-drop task management with real-time updates</p>
          </div>
          <div className="bg-slate-800/50 rounded-xl p-6 border border-slate-700">
            <div className="text-3xl mb-4">📅</div>
            <h3 className="text-xl font-semibold text-white mb-2">Calendar View</h3>
            <p className="text-slate-400">Visualize deadlines and milestones at a glance</p>
          </div>
          <div className="bg-slate-800/50 rounded-xl p-6 border border-slate-700">
            <div className="text-3xl mb-4">📊</div>
            <h3 className="text-xl font-semibold text-white mb-2">Analytics</h3>
            <p className="text-slate-400">Track progress with detailed reporting dashboards</p>
          </div>
        </div>

        {/* Tech Stack */}
        <div className="text-center text-slate-400 text-sm">
          Built with Next.js • Go • PostgreSQL • Redis
        </div>
      </main>
    </div>
  )
}
