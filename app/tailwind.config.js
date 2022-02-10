module.exports = {
  mode: "jit",
  future: {
    removeDeprecatedGapUtilities: true,
    purgeLayersByDefault: true,
  },
  content: ["./src/*.html", "./src/**/*.tsx"],
  theme: {
    extend: {
      spacing: {
        '704': '44rem',
      }
    },
  },
  variants: {},
  plugins: [],
};
