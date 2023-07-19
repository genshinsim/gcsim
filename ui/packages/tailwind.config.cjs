// eslint-disable-next-line @typescript-eslint/no-var-requires
const { Colors } = require("@blueprintjs/core");

module.exports = {
  future: {
    removeDeprecatedGapUtilities: true,
    purgeLayersByDefault: true,
  },
  content: [
    "./index.html",
    "./src/**/*.{html,js,ts,jsx,tsx}",
    "../ui/src/**/*.{js,ts,jsx,tsx}"
  ],
  theme: {
    screens: {
      xs: "400px",
      // 25rem
      sm: "640px",
      // => @media (min-width: 640px) { ... }

      md: "768px",
      // => @media (min-width: 768px) { ... }

      hd: "1000px",

      lg: "1024px",
      // => @media (min-width: 1024px) { ... }

      wide: "1160px",

      xl: "1280px",
      // => @media (min-width: 1280px) { ... }

      "2xl": "1536px",
      // => @media (min-width: 1536px) { ... }
    },
    extend: {
      spacing: { 320: "320px", 280: "280" },
      colors: {
        "bp-header-color": "#394b59",
        "bp-card-color": "#30404d",
        "bp-bg": "#293742",

        // https://blueprintjs.com/docs/#core/colors
        "bp4-black": Colors.BLACK,
        "bp4-dark-gray": {
          100: Colors.DARK_GRAY1,
          200: Colors.DARK_GRAY2,
          300: Colors.DARK_GRAY3,
          400: Colors.DARK_GRAY4,
          500: Colors.DARK_GRAY5,
        },
        "bp4-gray": {
          100: Colors.GRAY1,
          200: Colors.GRAY2,
          300: Colors.GRAY3,
          400: Colors.GRAY4,
          500: Colors.GRAY5,
        },
        "bp4-light-gray": {
          100: Colors.LIGHT_GRAY1,
          200: Colors.LIGHT_GRAY2,
          300: Colors.LIGHT_GRAY3,
          400: Colors.LIGHT_GRAY4,
          500: Colors.LIGHT_GRAY5,
        },

        "anemo": "#61DBBB",
        "geo": "#F8BA4E",
        "electro": "#B25DCD",
        "hydro": "#2F63D4",
        "pyro": "#BF2818",
        "cryo": "#77A2E6",
        "dendro": "#A5C83B",
      },
    },
  },
  variants: {},
  plugins: [],
};
