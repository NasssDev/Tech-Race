/** @type {import('tailwindcss').Config} */
module.exports = {
  content: ["./views/*.html"],
  theme: {
    extend: {
      fontFamily: {
        sans: ["Iosevka Aile Iaso", "sans-serif"],
        mono: ["Iosevka Curly Iaso", "monospace"],
        serif: ["Iosevka Etoile Iaso", "serif"],
        lexend: ["var(--font-family-lexend-deca)"],
        shoulder: ["var(--font-family-big-shoulder)"]
      },
      colors: {
        'bright-orange': "var(--Bright-orange)",
        'cyan': "var(--Dark-cyan)",
        'dark-cyan': "var(--Very-dark-cyan)",
        'transparent-white': "var(--Transparent-white)",
        'light-gray': "var(--Very-light-gray)"
      },
      spacing: {
        'second-sub-level': "var(--second-sub-level)", // 4px
        'first-sub-level': "var(--first-sub-level)",   // 8px
        'root-space': "var(--root-space)",             // 16px
        'first-level-space': "var(--first-level-space)", // 32px
        'second-level-space': "var(--second-level-space)", // 40px or 64px (reused variable name)
      }
    },
  },
  plugins: [],
};
