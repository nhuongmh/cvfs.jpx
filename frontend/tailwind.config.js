const {heroui} = require('@heroui/theme');
const {nextui} = require('@nextui-org/theme');
const defaultTheme = require("tailwindcss/defaultTheme");

/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./app/**/*.{js,ts,jsx,tsx}",
    "./pages/**/*.{js,ts,jsx,tsx}",
    "./components/**/*.{js,ts,jsx,tsx}",
    "./node_modules/@nextui-org/theme/dist/**/*.{js,ts,jsx,tsx}",
    "./node_modules/@heroui/theme/dist/components/(table|checkbox|form|spacer).js"
  ],
  theme: {
    fontFamily: {
      sans: ["var(--font-inter)", ...defaultTheme.fontFamily.sans],
    },
    extend: {
      dropShadow: {
        cta: ["0 10px 15px rgba(219, 227, 248, 0.2)"],
        blue: ["0 10px 15px rgba(59, 130, 246, 0.2)"],
      },
    },
  },
  darkMode: "class",
  plugins: [require("@tailwindcss/forms"),nextui(),heroui()],
};
