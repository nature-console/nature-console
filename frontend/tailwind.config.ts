import type { Config } from 'tailwindcss'

const config: Config = {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        'nature-bg': '#faf9f5',
        'nature-text': '#141414',
        'nature-green': '#bcd1ca',
        'nature-orange': '#d97756',
        'nature-purple': '#cbcadb',
        'nature-beige': '#f0eee6',
      },
    },
  },
  plugins: [],
}
export default config