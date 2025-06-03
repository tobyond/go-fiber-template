import esbuild from 'esbuild'
import fs from 'fs'
import path from 'path'

const isProduction = process.env.NODE_ENV === 'production'
const isWatch = process.argv.includes('--watch')

// Ensure output directory exists
const distDir = 'web/static/js'
if (!fs.existsSync(distDir)) {
  fs.mkdirSync(distDir, { recursive: true })
}

// Build configuration for JavaScript only
const buildConfig = {
  entryPoints: ['web/assets/js/app.js'],
  bundle: true,
  outdir: distDir,
  sourcemap: !isProduction,
  minify: isProduction,
  target: ['es2020'],
  format: 'esm',
  entryNames: isProduction ? '[name]-[hash]' : '[name]',
}

async function build() {
  try {
    console.log(isWatch ? '👀 Watching JavaScript files...' : '🔨 Building JavaScript...')
    
    if (isWatch) {
      const ctx = await esbuild.context(buildConfig)
      await ctx.watch()
      console.log('✅ JavaScript watcher started!')
      
      // Keep process alive
      process.on('SIGINT', async () => {
        console.log('\n👋 Stopping JavaScript watcher...')
        await ctx.dispose()
        process.exit(0)
      })
    } else {
      await esbuild.build(buildConfig)
      console.log('✅ JavaScript build complete!')
    }
    
  } catch (error) {
    console.error('❌ JavaScript build failed:', error)
    process.exit(1)
  }
}

build()
