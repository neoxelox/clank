module.exports = {
  content: ["./src/**/*.{html,js,svelte,ts}", "./node_modules/layerchart/**/*.{svelte,js}"],
  darkMode: "selector",
  theme: {
    container: {
      center: true,
      padding: "2rem",
      screens: {
        "2xl": "1400px",
      },
    },
    extend: {
      colors: {
        border: "hsl(var(--border) / <alpha-value>)",
        input: "hsl(var(--input) / <alpha-value>)",
        ring: "hsl(var(--ring) / <alpha-value>)",
        background: "hsl(var(--background) / <alpha-value>)",
        foreground: "hsl(var(--foreground) / <alpha-value>)",
        primary: {
          DEFAULT: "hsl(var(--primary) / <alpha-value>)",
          lighter: "hsl(var(--primary-lighter) / <alpha-value>)",
          foreground: "hsl(var(--primary-foreground) / <alpha-value>)",
        },
        secondary: {
          DEFAULT: "hsl(var(--secondary) / <alpha-value>)",
          darker: "hsl(var(--secondary-darker) / <alpha-value>)",
          foreground: "hsl(var(--secondary-foreground) / <alpha-value>)",
        },
        constructive: {
          DEFAULT: "hsl(var(--constructive) / <alpha-value>)",
          foreground: "hsl(var(--constructive-foreground) / <alpha-value>)",
        },
        destructive: {
          DEFAULT: "hsl(var(--destructive) / <alpha-value>)",
          foreground: "hsl(var(--destructive-foreground) / <alpha-value>)",
        },
        muted: {
          DEFAULT: "hsl(var(--muted) / <alpha-value>)",
          foreground: "hsl(var(--muted-foreground) / <alpha-value>)",
        },
        accent: {
          DEFAULT: "hsl(var(--accent) / <alpha-value>)",
          foreground: "hsl(var(--accent-foreground) / <alpha-value>)",
        },
        popover: {
          DEFAULT: "hsl(var(--popover) / <alpha-value>)",
          foreground: "hsl(var(--popover-foreground) / <alpha-value>)",
        },
        card: {
          DEFAULT: "hsl(var(--card) / <alpha-value>)",
          foreground: "hsl(var(--card-foreground) / <alpha-value>)",
        },
        surface: {
          DEFAULT: "hsl(var(--background) / <alpha-value>)",
          foreground: "hsl(var(--foreground) / <alpha-value>)",
          content: "hsl(var(--foreground) / <alpha-value>)",
          100: "hsl(var(--background) / <alpha-value>)",
          200: "hsl(var(---muted) / <alpha-value>)",
          300: "hsl(var(--border) / <alpha-value>)",
        },
      },
      borderRadius: {
        lg: "var(--radius)",
        md: "calc(var(--radius) - 2px)",
        sm: "calc(var(--radius) - 4px)",
      },
      fontFamily: {
        sans: "var(--font-geist-sans)",
        mono: "var(--font-geist-mono)",
        inter: "var(--font-inter-sans)",
        cal: "var(--font-cal-sans)",
      },
      animation: {
        "banner-appear": "banner-appear 5s ease-in",
        "marquee-lead": "marquee-lead 50s linear infinite",
        "marquee-rear": "marquee-rear 50s linear infinite",
        "gradient-first": "gradient-vertical 30s ease infinite",
        "gradient-second": "gradient-circle 20s reverse infinite",
        "gradient-third": "gradient-circle 40s linear infinite",
        "gradient-fourth": "gradient-horizontal 20s ease infinite",
        "gradient-fifth": "gradient-circle 20s ease infinite",
      },
      keyframes: {
        "banner-appear": {
          "0%": { opacity: 0 },
          "75%": { opacity: 0 },
          "100%": { opacity: 1 },
        },
        "marquee-lead": {
          "0%": { transform: "translateX(0%)" },
          "100%": { transform: "translateX(-100%)" },
        },
        "marquee-rear": {
          "0%": { transform: "translateX(100%)" },
          "100%": { transform: "translateX(0%)" },
        },
        "gradient-horizontal": {
          "0%": { transform: "translateX(-50%) translateY(-10%)" },
          "50%": { transform: "translateX(50%) translateY(10%)" },
          "100%": { transform: "translateX(-50%) translateY(-10%)" },
        },
        "gradient-circle": {
          "0%": { transform: "rotate(0deg)" },
          "50%": { transform: "rotate(180deg)" },
          "100%": { transform: "rotate(360deg)" },
        },
        "gradient-vertical": {
          "0%": { transform: "translateY(-50%)" },
          "50%": { transform: "translateY(50%)" },
          "100%": { transform: "translateY(-50%)" },
        },
      },
    },
  },
  variants: {
    extend: {},
  },
  safelist: ["dark"],
  plugins: [
    require("@tailwindcss/aspect-ratio"),
    require("@tailwindcss/container-queries"),
    require("@tailwindcss/forms"),
    require("@tailwindcss/typography"),
  ],
};
