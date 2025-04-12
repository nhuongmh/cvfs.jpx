export type SiteConfig = typeof siteConfig;

export const siteConfig = {
    name: ".langfi",
    description: "langfi: Accelerate your language learning",
    navItems: [
        {
            name: "Stat",
            href: "/stat",
        },
        {
            name: "Process",
            href: "/process",
        },

        {
            name: "Learn",
            href: "/learn",
            external: false,
        },
    ],
    navMenuItems: [
        {
            label: "Profile",
            href: "/profile",
        },
        {
            label: "Dashboard",
            href: "/dashboard",
        },
        {
            label: "Projects",
            href: "/projects",
        },
        {
            label: "Team",
            href: "/team",
        },
        {
            label: "Calendar",
            href: "/calendar",
        },
        {
            label: "Settings",
            href: "/settings",
        },
        {
            label: "Help & Feedback",
            href: "/help-feedback",
        },
        {
            label: "Logout",
            href: "/logout",
        },
    ],
    links: {
        github: "https://github.com/nhuongmh/cvfs.jpx",
        docs: "https://nextui.org",
        googlechat: "https://chat.google.com/room/AAAA-JatYZI?cls=7",
    },
    server_url_prefix: `http://${process.env.NEXT_PUBLIC_SERVER_URL}/public/api/v1`,
    private_url_prefix: `http://${process.env.NEXT_PUBLIC_SERVER_URL}/private/api/v1`
};