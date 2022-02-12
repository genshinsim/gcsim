module.exports = {
  mode: "jit",
  future: {
    removeDeprecatedGapUtilities: true,
    purgeLayersByDefault: true,
  },
  content: ["./src/*.html", "./src/**/*.tsx"],
  theme: {
    screens: {
      xs: "400px",
      // 25rem
      sm: "640px",
      // => @media (min-width: 640px) { ... }

      md: "768px",
      // => @media (min-width: 768px) { ... }

      wide: "1160px",

      lg: "1024px",
      // => @media (min-width: 1024px) { ... }

      xl: "1280px",
      // => @media (min-width: 1280px) { ... }

      "2xl": "1536px",
      // => @media (min-width: 1536px) { ... }
    },
    extend: {
      spacing: { 320: "320px" },
      colors: {
        "bp-header-color": "#394b59",
        "bp-card-color": "#30404d",
        "bp-bg": "#293742",
      },
    },
  },
  variants: {},
  plugins: [],
};
