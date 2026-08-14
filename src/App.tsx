import { Outlet, Router } from "@/app/Router"
import { RootRoute } from "@/routes/Root"

export function App() {
  return (
    <Router>
      <RootRoute><Outlet /></RootRoute>
    </Router>
  )
}

export default App
