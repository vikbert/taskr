import { defineConfig, HeadConfig } from 'vitepress';
import githubLinksPlugin from './plugins/github-links';
import { readFileSync } from 'fs';
import { resolve } from 'path';
import { tabsMarkdownPlugin } from 'vitepress-plugin-tabs';
import {
  groupIconMdPlugin,
  groupIconVitePlugin,
  localIconLoader
} from 'vitepress-plugin-group-icons';
import { taskDescription, taskName, ogUrl, ogImage } from './meta.ts';
import { fileURLToPath, URL } from 'node:url';
import llmstxt, { copyOrDownloadAsMarkdownButtons } from 'vitepress-plugin-llms';

const version = readFileSync(
  resolve(__dirname, '../../internal/version/version.txt'),
  'utf8'
).trim();

const urlVersion =
  process.env.NODE_ENV === 'development'
    ? {
        current: 'https://taskr-io.vercel.app/',
        next: 'http://localhost:3002/'
      }
    : {
        current: 'https://taskr-io.vercel.app/',
        next: 'https://next.taskr-io.vercel.app/'
      };

// https://vitepress.dev/reference/site-config
export default defineConfig({
  title: taskName,
  description: taskDescription,
  lang: 'en-US',
  head: [
    [
      'link',
      {
        rel: 'icon',
        type: 'image/x-icon',
        href: '/img/favicon.ico',
        sizes: '48x48'
      }
    ],
    [
      'link',
      {
        rel: 'icon',
        sizes: 'any',
        type: 'image/png',
        href: '/img/logo.png'
      }
    ],
    // Open Graph
    ['meta', { property: 'og:type', content: 'website' }],
    ['meta', { property: 'og:site_name', content: taskName }],
    ['meta', { property: 'og:title', content: taskName }],
    ['meta', { property: 'og:description', content: taskDescription }],
    ['meta', { property: 'og:image', content: ogImage }],
    ['meta', { property: 'og:url', content: ogUrl }],
    // Twitter Card
    [
      'meta',
      {
        name: 'keywords',
        content:
          'task runner, build tool, taskfile, yaml build tool, go task runner, make alternative, cross-platform build tool, makefile alternative, automation tool, ci cd pipeline, developer productivity, build automation, command line tool, go binary, yaml configuration'
      }
    ],
    [
      "script",
      {
        defer: "",
        src: "https://u.taskr-io.vercel.app/script.js",
        "data-website-id": "084030b0-0e3f-4891-8d2a-0c12c40f5933"
      }
    ]
  ],
  transformHead({ pageData }) {
    const head: HeadConfig[] = []

    // Canonical URL dynamique
    const canonicalUrl = `https://taskr-io.vercel.app/${pageData.relativePath
      .replace(/\.md$/, '')
      .replace(/index$/, '')}`
    head.push(['link', { rel: 'canonical', href: canonicalUrl }])

    // Noindex pour 404
    if (pageData.relativePath === '404.md') {
      head.push(['meta', { name: 'robots', content: 'noindex, nofollow' }])
    }

    return head
  },
  srcDir: 'src',
  cleanUrls: true,
  markdown: {
    config: (md) => {
      md.use(githubLinksPlugin, {
        baseUrl: 'https://github.com',
        repo: 'vikbert/taskr'
      });
      md.use(tabsMarkdownPlugin);
      md.use(groupIconMdPlugin);
      md.use(copyOrDownloadAsMarkdownButtons);
    }
  },
  vite: {
    plugins: [
      llmstxt({
        ignoreFiles: [
          'index.md',
          'team.md',
          'donate.md',
          'docs/styleguide.md',
          'docs/contributing.md',
          'docs/releasing.md',
          'docs/changelog.md',
          'blog/*'
        ]
      }),
      groupIconVitePlugin({
        customIcon: {
          '.taskrc.yml': localIconLoader(
            import.meta.url,
            './theme/icons/task.svg'
          ),
          'Taskfile.yml': localIconLoader(
            import.meta.url,
            './theme/icons/task.svg'
          )
        }
      })
    ],
    resolve: {
      alias: [
        {
          find: /^.*\/VPTeamMembersItem\.vue$/,
          replacement: fileURLToPath(
            new URL('./components/VPTeamMembersItem.vue', import.meta.url)
          )
        }
      ]
    }
  },

  themeConfig: {
    logo: '/img/logo.png',
    carbonAds: {
      code: 'CESI65QJ',
      placement: 'taskfiledev'
    },
    search: {
      provider: 'algolia',
      options: {
        appId: '7IZIJ13AI7',
        apiKey: '34b64ae4fc8d9da43d9a13d9710aaddc',
        indexName: 'taskfile'
      }
    },
    nav: [
      { text: 'Home', link: '/' },
      {
        text: 'Docs',
        link: '/docs/guide',
        activeMatch: '^/docs'
      }
    ],

    sidebar: {
      '/': [
        {
          text: 'Installation',
          link: '/docs/installation'
        },
        {
          text: 'Getting Started',
          link: '/docs/getting-started'
        },
        {
          text: 'Guide',
          link: '/docs/guide'
        },
        {
          text: 'Reference',
          collapsed: true,
          items: [
            {
              text: 'Taskfile Schema',
              link: '/docs/reference/schema'
            },
            {
              text: 'Environment',
              link: '/docs/reference/environment'
            },
            {
              text: 'Configuration',
              link: '/docs/reference/config'
            },
            {
              text: 'CLI',
              link: '/docs/reference/cli'
            },
            {
              text: 'Templating',
              link: '/docs/reference/templating'
            },
            {
              text: 'Package API',
              link: '/docs/reference/package'
            }
          ]
        },
        {
          text: 'Taskfile Versions',
          link: '/docs/taskfile-versions'
        },
        {
          text: 'Integrations',
          link: '/docs/integrations'
        },
        {
          text: 'Community',
          link: '/docs/community'
        },
        {
          text: 'Style Guide',
          link: '/docs/styleguide'
        },
        {
          text: 'Contributing',
          link: '/docs/contributing'
        },
        {
          text: 'Releasing',
          link: '/docs/releasing'
        },
        {
          text: 'Changelog',
          link: '/docs/changelog'
        },
        {
          text: 'FAQ',
          link: '/docs/faq'
        }
      ],
    },

    socialLinks: [
      { icon: 'github', link: 'https://github.com/vikbert/taskr' },
    ],

    footer: {
      message:
        'Built with <a target="_blank" href="https://www.vercel.com">Vercel</a>'
    }
  },
  sitemap: {
    hostname: 'https://taskr-io.vercel.app',
    transformItems: (items) => {
      return items.map((item) => ({
        ...item,
        lastmod: new Date().toISOString()
      }));
    }
  }
});
