import { Boxes } from "lucide-react"

export const apps = [
  { name: "Telegram", desc: "Connect with Telegram for real-time communication.", connected: false },
  { name: "Notion", desc: "Effortlessly sync Notion pages for seamless collaboration.", connected: true },
  { name: "Figma", desc: "View and collaborate on Figma designs in one place.", connected: true },
  { name: "Trello", desc: "Sync Trello cards for streamlined project management.", connected: false },
  { name: "Slack", desc: "Integrate Slack for efficient team communication.", connected: false },
  { name: "Zoom", desc: "Host Zoom meetings directly from the dashboard.", connected: true },
  { name: "Stripe", desc: "Easily manage Stripe transactions and payments.", connected: false },
  { name: "Gmail", desc: "Access and manage Gmail messages effortlessly.", connected: true },
  { name: "GitHub", desc: "Streamline code management with GitHub integration.", connected: false },
  { name: "Docker", desc: "Effortlessly manage Docker containers on your dashboard.", connected: false },
].map((app) => ({ ...app, logo: Boxes }))
