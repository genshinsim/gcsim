// @ts-check
// Note: type annotations allow type checking and IDEs autocompletion
// const math = require('remark-math');
// const katex = require('rehype-katex');
const lightCodeTheme = require('prism-react-renderer/themes/github');
const darkCodeTheme = require('prism-react-renderer/themes/dracula');

async function createConfig() {
  const katex = (await import('rehype-katex')).default;
  const math = (await import('remark-math')).default;

  /** @type {import('@docusaurus/types').Config} */
  const config = {
    title: "gcsim Docs",
    tagline: "gcsim - simulation impact",
    url: "https://docs.gcsim.app",
    baseUrl: "/",
    onBrokenLinks: 'throw',
    onBrokenMarkdownLinks: 'warn',
    scripts: [
      {
        src: "https://static.cloudflareinsights.com/beacon.min.js",
        defer: true,
        "data-cf-beacon": '{"token": "2f8a17efd29b479f9dc27e09aae7ccb5"}',
      },
    ],

    // Even if you don't use internalization, you can use this field to set useful
    // metadata like html lang. For example, if your site is Chinese, you may want
    // to replace "en" with "zh-Hans".
    i18n: {
      defaultLocale: 'en',
      locales: ['en'],
    },

    presets: [
      [
        '@docusaurus/preset-classic',
        // /** @type {import('@docusaurus/preset-classic').Options} */
        ({
          docs: {
            routeBasePath: "/",
            sidebarPath: require.resolve('./sidebars.js'),
            editUrl: "https://github.com/genshinsim/gcsim/blob/main/ui/packages/docs",
            remarkPlugins: [math],
            rehypePlugins: [katex],
          },
          blog: false,
          theme: {
            customCss: require.resolve('./src/css/custom.css'),
          },
        }),
      ],
    ],
    stylesheets: [
      {
        href: 'https://cdn.jsdelivr.net/npm/katex@0.15.2/dist/katex.min.css',
        type: 'text/css',
        integrity:
          'sha384-MlJdn/WNKDGXveldHDdyRP1R4CTHr3FeuDNfhsLPYrq2t0UBkUdK2jyTnXPEK1NQ',
        crossorigin: 'anonymous',
      },
    ],

    themeConfig:
      /** @type {import('@docusaurus/preset-classic').ThemeConfig} */
      ({
        // Replace with your project's social card
        navbar: {
          title: "gcsim Docs",
          items: [
            {
              label: "Discord",
              href: "https://discord.gg/m7jvjdxx7q",
              position: "right",
            },
            {
              href: "https://github.com/genshinsim/gsim",
              label: "GitHub",
              position: "right",
            },
          ],
        },
        footer: {
          style: "dark",
        },
        prism: {
          theme: lightCodeTheme,
          darkTheme: darkCodeTheme,
        },
        colorMode: {
          defaultMode: "dark",
          disableSwitch: false,
          respectPrefersColorScheme: false,
        },
        algolia: {
          // The application ID provided by Algolia
          appId: 'FQ95W3KA6U',
    
          // Public API key: it is safe to commit it
          apiKey: '3d7fffa98beeefe3652e892d29937cf9',
    
          indexName: 'gcsim',
    
          // Optional: see doc section below
          contextualSearch: true,
    
        },
      }),

    plugins: [
      async function myPlugin(context, options) {
        return {
          name: "docusaurus-tailwindcss",
          configurePostCss(postcssOptions) {
            // Appends TailwindCSS and AutoPrefixer.
            postcssOptions.plugins.push(require("tailwindcss"));
            postcssOptions.plugins.push(require("autoprefixer"));
            return postcssOptions;
          },
        };
      },
    ],

  };
  return config
}


module.exports = createConfig;
