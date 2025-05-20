#!/usr/bin/env node

import fs from 'fs';
import path from 'path';
import { fileURLToPath } from 'url';

// Get __dirname equivalent in ESM
const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

async function findSvelteFiles(dir) {
  const files = [];
  const entries = await fs.promises.readdir(dir, { withFileTypes: true });

  for (const entry of entries) {
    const fullPath = path.join(dir, entry.name);
    
    if (entry.isDirectory() && entry.name !== 'node_modules' && !entry.name.startsWith('.')) {
      const subFiles = await findSvelteFiles(fullPath);
      files.push(...subFiles);
    } else if (entry.name.endsWith('.svelte')) {
      files.push(fullPath);
    }
  }

  return files;
}

function convertToRunes(content) {
  // Replace reactive declarations with $derived
  content = content.replace(/\$:\s+([a-zA-Z0-9_$]+)\s*=\s*(.*?);/g, (match, varName, expression) => {
    return `const ${varName} = $derived(${expression});`;
  });

  // Replace reactive statements with $effect
  content = content.replace(/\$:\s+\{([\s\S]*?)\}/g, (match, code) => {
    return `$effect(() => {${code}});`;
  });

  // Replace reactive statements without braces with $effect
  content = content.replace(/\$:\s+(?!\{)((?!const)[^;]*);/g, (match, code) => {
    if (!code.includes('=')) { // Only if it's not an assignment
      return `$effect(() => { ${code}; });`;
    }
    return match;
  });

  // Replace store subscriptions with $ prefix but not affect existing $ prefixed variables
  // This is complex and might require more careful handling

  return content;
}

async function convertFile(filePath) {
  console.log(`Converting ${filePath}`);
  let content = await fs.promises.readFile(filePath, 'utf8');
  content = convertToRunes(content);
  await fs.promises.writeFile(filePath, content, 'utf8');
}

async function main() {
  const baseDir = path.join(__dirname, 'src');
  const svelteFiles = await findSvelteFiles(baseDir);
  console.log(`Found ${svelteFiles.length} Svelte files`);
  
  for (const file of svelteFiles) {
    await convertFile(file);
  }
  
  console.log('Done!');
}

main().catch(console.error); 