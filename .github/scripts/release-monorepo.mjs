#!/usr/bin/env node

/**
 * Chia Florist — Automated Monorepo Tag & Release Engine
 *
 * Direct Git tag & GitHub Release automation for independent monorepo packages.
 * - Zero bot-authored commits
 * - Zero automated release pull requests / release branches
 * - Independent versioning per package based on directory path changes
 * - Semantic versioning with Conventional Commits parsing
 * - Pre-major minor bump support (0.x.y)
 */

import { execSync } from 'node:child_process';
import { readFileSync, existsSync } from 'node:fs';
import { resolve, dirname } from 'node:path';
import { fileURLToPath } from 'node:url';

const __dirname = dirname(fileURLToPath(import.meta.url));

// --- CLI Args & Options --- //
const args = process.argv.slice(2);
const isDryRun = args.includes('--dry-run');
const configIndex = args.indexOf('--config');
const configPath = configIndex !== -1 && args[configIndex + 1]
  ? resolve(process.cwd(), args[configIndex + 1])
  : resolve(__dirname, '../release-config.json');

const token = process.env.GH_TOKEN || process.env.GITHUB_TOKEN || process.env.CHIA_CURRENT_ARC;
const repository = process.env.GITHUB_REPOSITORY; // e.g. "aether-gear/chia.florist"
const targetSha = process.env.GITHUB_SHA || getGitHeadSha();

function getGitHeadSha() {
  try {
    return execSync('git rev-parse HEAD', { encoding: 'utf-8' }).trim();
  } catch {
    return 'HEAD';
  }
}

function runGit(cmd) {
  return execSync(cmd, { encoding: 'utf-8', stdio: ['pipe', 'pipe', 'pipe'] }).trim();
}

// --- Semver Utilities --- //
function parseSemver(str) {
  const match = str.match(/^v?(\d+)\.(\d+)\.(\d+)(?:-([0-9A-Za-z.-]+))?$/);
  if (!match) return null;
  return {
    major: parseInt(match[1], 10),
    minor: parseInt(match[2], 10),
    patch: parseInt(match[3], 10),
    prerelease: match[4] || null,
    raw: str
  };
}

function compareSemver(a, b) {
  if (a.major !== b.major) return a.major - b.major;
  if (a.minor !== b.minor) return a.minor - b.minor;
  if (a.patch !== b.patch) return a.patch - b.patch;
  if (a.prerelease && !b.prerelease) return -1;
  if (!a.prerelease && b.prerelease) return 1;
  return 0;
}

function bumpVersion(current, bumpType, bumpMinorPreMajor = true) {
  const { major, minor, patch } = current;

  if (bumpType === 'major') {
    if (bumpMinorPreMajor && major === 0) {
      return { major: 0, minor: minor + 1, patch: 0 };
    }
    return { major: major + 1, minor: 0, patch: 0 };
  }

  if (bumpType === 'minor') {
    return { major, minor: minor + 1, patch: 0 };
  }

  if (bumpType === 'patch') {
    return { major, minor, patch: patch + 1 };
  }

  return current;
}

function formatSemver(v) {
  return `${v.major}.${v.minor}.${v.patch}${v.prerelease ? `-${v.prerelease}` : ''}`;
}

// --- Conventional Commits Parser --- //
const COMMIT_REGEX = /^([a-zA-Z]+)(?:\(([^)]+)\))?(!)?:\s*(.+)$/;

function parseCommit(rawHash, rawSubject, rawBody) {
  const hash = rawHash.trim();
  const subject = rawSubject.trim();
  const body = rawBody.trim();

  let type = 'other';
  let scope = null;
  let isBreaking = false;
  let description = subject;

  const match = subject.match(COMMIT_REGEX);
  if (match) {
    type = match[1].toLowerCase();
    scope = match[2] || null;
    if (match[3] === '!') isBreaking = true;
    description = match[4].trim();
  }

  // Check body for breaking change footer
  if (
    body.includes('BREAKING CHANGE:') ||
    body.includes('BREAKING-CHANGE:') ||
    subject.includes('BREAKING CHANGE:')
  ) {
    isBreaking = true;
  }

  let bump = 'none';
  if (isBreaking) {
    bump = 'major';
  } else if (type === 'feat') {
    bump = 'minor';
  } else if (type === 'fix' || type === 'perf') {
    bump = 'patch';
  }

  return {
    hash,
    shortHash: hash.slice(0, 8),
    type,
    scope,
    description,
    subject,
    body,
    isBreaking,
    bump
  };
}

// --- Tag Discovery per Package --- //
function findLatestPackageTag(tagPrefix) {
  try {
    const rawTags = runGit(`git tag -l "${tagPrefix}*"`);
    if (!rawTags) return null;

    const tags = rawTags
      .split('\n')
      .map(t => t.trim())
      .filter(Boolean)
      .map(tag => {
        const versionStr = tag.slice(tagPrefix.length);
        const semver = parseSemver(versionStr);
        return semver ? { tag, semver } : null;
      })
      .filter(Boolean);

    if (tags.length === 0) return null;

    tags.sort((a, b) => compareSemver(a.semver, b.semver));
    return tags[tags.length - 1];
  } catch (err) {
    console.warn(`[warn] Could not list tags for ${tagPrefix}:`, err.message);
    return null;
  }
}

function isTagAlreadyAtSha(tagName, sha) {
  try {
    const tagsAtSha = runGit(`git tag --points-at ${sha}`)
      .split('\n')
      .map(t => t.trim())
      .filter(Boolean);
    return tagsAtSha.includes(tagName);
  } catch {
    return false;
  }
}

// --- Get Commits Touching Package Path --- //
function getCommitsForPackage(pkgPath, latestTag, target) {
  try {
    const range = latestTag ? `${latestTag.tag}..${target}` : target;
    // Format: %H (hash) %x1f %s (subject) %x1f %b (body) %x1e (record separator)
    const logCmd = `git log ${range} --format="%H%x1f%s%x1f%b%x1e" -- "${pkgPath}"`;
    const output = runGit(logCmd);

    if (!output) return [];

    return output
      .split('\x1e')
      .map(entry => entry.trim())
      .filter(Boolean)
      .map(entry => {
        const [hash = '', subject = '', body = ''] = entry.split('\x1f');
        return parseCommit(hash, subject, body);
      });
  } catch (err) {
    console.warn(`[warn] Error getting git log for ${pkgPath}:`, err.message);
    return [];
  }
}

// --- Release Notes Generator --- //
function generateReleaseNotes(pkgName, newTag, prevTag, commits, repoUrl) {
  const breaking = commits.filter(c => c.isBreaking);
  const features = commits.filter(c => c.type === 'feat' && !c.isBreaking);
  const fixes = commits.filter(c => (c.type === 'fix' || c.type === 'perf') && !c.isBreaking);
  const others = commits.filter(c => !c.isBreaking && c.type !== 'feat' && c.type !== 'fix' && c.type !== 'perf');

  const lines = [];
  lines.push(`## ${newTag}`);
  lines.push('');

  const formatItem = (c) => {
    const scopePrefix = c.scope ? `**${c.scope}:** ` : '';
    const link = repoUrl ? ` ([${c.shortHash}](${repoUrl}/commit/${c.hash}))` : ` (${c.shortHash})`;
    return `* ${scopePrefix}${c.description}${link}`;
  };

  if (breaking.length > 0) {
    lines.push('### 💥 Breaking Changes');
    breaking.forEach(c => lines.push(formatItem(c)));
    lines.push('');
  }

  if (features.length > 0) {
    lines.push('### ✨ Features');
    features.forEach(c => lines.push(formatItem(c)));
    lines.push('');
  }

  if (fixes.length > 0) {
    lines.push('### 🐛 Bug Fixes & Performance');
    fixes.forEach(c => lines.push(formatItem(c)));
    lines.push('');
  }

  if (others.length > 0 && (breaking.length > 0 || features.length > 0 || fixes.length > 0)) {
    lines.push('### 🛠️ Maintenance & Refactoring');
    others.forEach(c => lines.push(formatItem(c)));
    lines.push('');
  }

  if (repoUrl && prevTag) {
    lines.push(`**Full Changelog**: ${repoUrl}/compare/${prevTag}...${newTag}`);
  }

  return lines.join('\n').trim();
}

// --- GitHub Release API Caller --- //
async function createGitHubRelease(newTag, targetCommitSha, releaseNotes) {
  if (!repository) {
    throw new Error('GITHUB_REPOSITORY env var is missing');
  }
  if (!token) {
    throw new Error('GH_TOKEN or GITHUB_TOKEN env var is missing');
  }

  const url = `https://api.github.com/repos/${repository}/releases`;
  const payload = {
    tag_name: newTag,
    target_commitish: targetCommitSha,
    name: newTag,
    body: releaseNotes,
    draft: false,
    prerelease: false
  };

  const res = await fetch(url, {
    method: 'POST',
    headers: {
      'Accept': 'application/vnd.github+json',
      'Authorization': `Bearer ${token}`,
      'X-GitHub-Api-Version': '2022-11-28',
      'Content-Type': 'application/json'
    },
    body: JSON.stringify(payload)
  });

  if (!res.ok) {
    const errorBody = await res.text();
    if (res.status === 422 && errorBody.includes('already_exists')) {
      console.log(`ℹ️ Release or tag ${newTag} already exists on GitHub. Skipping.`);
      return { html_url: `https://github.com/${repository}/releases/tag/${newTag}` };
    }
    throw new Error(`GitHub API release creation failed (${res.status}): ${errorBody}`);
  }

  return await res.json();
}

// --- Main Engine --- //
async function main() {
  console.log('🌸 Chia Florist Monorepo Release Engine');
  console.log('========================================');
  if (isDryRun) {
    console.log('🔍 MODE: DRY-RUN (No tags or releases will be created)');
  }
  console.log(`📌 Target SHA : ${targetSha}`);
  console.log(`📦 Repository : ${repository || 'local (no GITHUB_REPOSITORY set)'}`);

  if (!existsSync(configPath)) {
    console.error(`❌ Config file not found at: ${configPath}`);
    process.exit(1);
  }

  const config = JSON.parse(readFileSync(configPath, 'utf-8'));
  const bumpMinorPreMajor = config.bumpMinorPreMajor ?? true;
  const packages = config.packages || [];

  const repoUrl = repository ? `https://github.com/${repository}` : '';
  const plannedReleases = [];

  for (const pkg of packages) {
    console.log(`\n────────────────────────────────────────`);
    console.log(`📦 Inspecting package: ${pkg.name} (${pkg.path})`);

    const latestTag = findLatestPackageTag(pkg.tagPrefix);
    const currentVersion = latestTag
      ? latestTag.semver
      : parseSemver(pkg.initialVersion || '0.0.0');

    console.log(`   Current Version : ${latestTag ? `${latestTag.tag} (${formatSemver(currentVersion)})` : `None (Baseline: ${formatSemver(currentVersion)})`}`);

    const commits = getCommitsForPackage(pkg.path, latestTag, targetSha);
    console.log(`   Commits touching '${pkg.path}': ${commits.length}`);

    if (commits.length === 0) {
      console.log(`   ↳ No changes. Skipping.`);
      continue;
    }

    // Determine highest bump
    let highestBump = 'none';
    for (const c of commits) {
      if (c.bump === 'major') {
        highestBump = 'major';
        break;
      } else if (c.bump === 'minor') {
        highestBump = 'minor';
      } else if (c.bump === 'patch' && highestBump === 'none') {
        highestBump = 'patch';
      }
    }

    if (highestBump === 'none') {
      console.log(`   ↳ Found ${commits.length} commit(s), but none trigger a version bump (only docs/chore/refactor/etc.). Skipping.`);
      continue;
    }

    const nextVersion = bumpVersion(currentVersion, highestBump, bumpMinorPreMajor);
    const nextSemverStr = formatSemver(nextVersion);
    const newTag = `${pkg.tagPrefix}${nextSemverStr}`;

    console.log(`   🚀 Bump Type     : ${highestBump.toUpperCase()}${bumpMinorPreMajor && currentVersion.major === 0 && highestBump === 'major' ? ' (pre-major minor bump)' : ''}`);
    console.log(`   ✨ New Version   : ${newTag}`);

    if (isTagAlreadyAtSha(newTag, targetSha)) {
      console.log(`   ↳ Tag ${newTag} is already pointing to target SHA. Skipping.`);
      continue;
    }

    const notes = generateReleaseNotes(
      pkg.name,
      newTag,
      latestTag ? latestTag.tag : null,
      commits,
      repoUrl
    );

    plannedReleases.push({
      pkg,
      newTag,
      prevTag: latestTag ? latestTag.tag : null,
      nextVersion: nextSemverStr,
      notes
    });
  }

  console.log(`\n========================================`);
  console.log(`📋 Summary: ${plannedReleases.length} package(s) to release`);

  if (plannedReleases.length === 0) {
    console.log('✅ Nothing to release. Exiting cleanly.');
    return;
  }

  for (const release of plannedReleases) {
    console.log(`\n----------------------------------------`);
    console.log(`Release: ${release.newTag}`);
    console.log(`Release Notes:\n${release.notes}\n`);

    if (isDryRun) {
      console.log(`[dry-run] Would create GitHub Release & tag: ${release.newTag} at SHA ${targetSha}`);
      continue;
    }

    try {
      console.log(`Creating GitHub Release & tag: ${release.newTag}...`);
      const res = await createGitHubRelease(release.newTag, targetSha, release.notes);
      console.log(`🎉 Successfully created release: ${res.html_url || release.newTag}`);
    } catch (err) {
      console.error(`❌ Failed to create release ${release.newTag}:`, err.message);
      process.exit(1);
    }
  }

  console.log(`\n✨ All releases completed successfully!`);
}

main().catch(err => {
  console.error('❌ Fatal error in release engine:', err);
  process.exit(1);
});
