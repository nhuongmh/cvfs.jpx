export type SiteConfig = typeof siteConfig;

export const siteConfig = {
    name: ".langfi",
    description: "langfi: Accelerate your language learning",
    navItems: [
        {
            label: "Home",
            href: "/",
        },
        {
            label: "Docs",
            href: "/docs",
        },
        {
            label: "Blog",
            href: "/blog",
        },
        {
            label: "Stat",
            href: "/share",
        },
        {
            label: "JPX",
            href: "/unseal",
        },

        {
            label: "Eng",
            href: "/deploy",
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
        github: "https://github.com/nextui-org/nextui",
        twitter: "https://twitter.com/getnextui",
        docs: "https://nextui.org",
        discord: "https://discord.gg/9b6yyZKmH4",
        sponsor: "https://patreon.com/jrgarciadev",
    },
    server_url_prefix: `http://${process.env.NEXT_PUBLIC_SERVER_URL}/public/api/v1`,
};