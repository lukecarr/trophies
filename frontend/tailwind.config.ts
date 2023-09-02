import type { Config } from "tailwindcss";

export default {
    content: [
        "./index.html",
        "./src/**/*.tsx",
    ],
    theme: {
        extend: {
            container: {
                center: true,
            },
        },
    },
    plugins: [],
} satisfies Config
