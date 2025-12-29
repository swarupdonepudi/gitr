import React, { useState } from 'react';
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Terminal, Copy, Check, Github, BookOpen, Star, GitBranch, Zap, Shield } from 'lucide-react';

export default function Hero() {
  const [copied, setCopied] = useState(false);

  const installCommand = 'brew install swarupdonepudi/tap/gitr';

  const copyToClipboard = () => {
    navigator.clipboard.writeText(installCommand);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="relative min-h-screen flex items-center">
      <div className="max-w-7xl mx-auto px-6 py-20">
        <div className="text-center max-w-4xl mx-auto">
          {/* Badges */}
          <div className="flex justify-center gap-3 mb-8 flex-wrap">
            <Badge variant="secondary" className="bg-gray-800/50 text-gray-300 border-gray-700">
              <Shield className="w-3 h-3 mr-1" />
              Apache-2.0
            </Badge>
            <Badge variant="secondary" className="bg-gray-800/50 text-gray-300 border-gray-700">
              <Terminal className="w-3 h-3 mr-1" />
              Go CLI
            </Badge>
            <Badge variant="secondary" className="bg-gray-800/50 text-gray-300 border-gray-700">
              <GitBranch className="w-3 h-3 mr-1" />
              GitHub/GitLab/Bitbucket
            </Badge>
            <Badge variant="secondary" className="bg-gray-800/50 text-gray-300 border-gray-700">
              <Zap className="w-3 h-3 mr-1" />
              SSH/HTTPS
            </Badge>
          </div>

          {/* Main Headline */}
          <h1 className="text-5xl md:text-7xl font-bold mb-6 bg-gradient-to-r from-white via-blue-100 to-purple-200 bg-clip-text text-transparent leading-tight">
            Clone smarter.
            <br />
            Jump to web faster.
          </h1>

          {/* Subheadline */}
          <p className="text-xl md:text-2xl text-gray-400 mb-12 max-w-3xl mx-auto leading-relaxed">
            <span className="text-blue-400 font-semibold">gitr</span> turns any repo URL into a{' '}
            <span className="text-purple-400">deterministic local path</span> and a{' '}
            <span className="text-green-400">one-shot deep link</span> in your browser.{' '}
            Less context switching, zero guesswork.
          </p>

          {/* Install Command */}
          <div className="bg-gray-900/80 backdrop-blur-xl border border-gray-700 rounded-2xl p-6 mb-8 max-w-4xl mx-auto">
            <div className="flex items-center justify-between mb-4">
              <div className="flex items-center gap-2">
                <div className="w-3 h-3 bg-red-500 rounded-full"></div>
                <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
                <div className="w-3 h-3 bg-green-500 rounded-full"></div>
                <span className="text-gray-500 text-sm ml-2">Terminal</span>
              </div>
              <Button
                variant="ghost"
                size="sm"
                onClick={copyToClipboard}
                className="text-gray-400 hover:text-white"
              >
                {copied ? (
                  <Check className="w-4 h-4 mr-2 text-green-400" />
                ) : (
                  <Copy className="w-4 h-4 mr-2" />
                )}
                {copied ? 'Copied!' : 'Copy'}
              </Button>
            </div>
            <code className="text-blue-400 text-sm md:text-base font-mono block">
              $ {installCommand}
            </code>
          </div>

          {/* CTA Buttons */}
          <div className="flex flex-col sm:flex-row items-center justify-center gap-4 mb-12">
            <Button 
              size="lg" 
              className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700 px-8 py-3 text-lg"
              onClick={() => document.getElementById('quickstart').scrollIntoView({ behavior: 'smooth' })}
            >
              Try an Example
            </Button>
            <Button 
              variant="outline" 
              size="lg" 
              className="border-gray-600 text-gray-300 hover:bg-gray-800 px-8 py-3 text-lg"
              asChild
            >
              <a href="https://github.com/swarupdonepudi/gitr" target="_blank" rel="noopener noreferrer">
                <Github className="w-5 h-5 mr-2" />
                View on GitHub
              </a>
            </Button>
          </div>

          {/* Demo Preview */}
          <div className="bg-gradient-to-r from-gray-900/50 to-gray-800/50 backdrop-blur-xl border border-gray-700 rounded-3xl p-8 max-w-5xl mx-auto">
            <div className="grid md:grid-cols-2 gap-8 items-center">
              <div className="space-y-4">
                <div className="bg-black/50 rounded-xl p-4 border border-gray-800">
                  <code className="text-green-400 text-sm font-mono">
                    $ gitr clone https://github.com/owner/repo
                  </code>
                  <div className="text-gray-500 text-xs mt-2 font-mono">
                    → Cloned to ~/scm/github.com/owner/repo
                  </div>
                </div>
                <div className="bg-black/50 rounded-xl p-4 border border-gray-800">
                  <code className="text-blue-400 text-sm font-mono">
                    $ gitr web
                  </code>
                  <div className="text-gray-500 text-xs mt-2 font-mono">
                    → Opens repo in browser
                  </div>
                </div>
              </div>
              <div className="space-y-3">
                <div className="text-left">
                  <h3 className="text-lg font-semibold text-white mb-2">What gitr does:</h3>
                  <ul className="space-y-2 text-gray-400">
                    <li className="flex items-center gap-2">
                      <div className="w-1.5 h-1.5 bg-blue-400 rounded-full"></div>
                      Parses repo URL structure
                    </li>
                    <li className="flex items-center gap-2">
                      <div className="w-1.5 h-1.5 bg-purple-400 rounded-full"></div>
                      Creates deterministic paths
                    </li>
                    <li className="flex items-center gap-2">
                      <div className="w-1.5 h-1.5 bg-green-400 rounded-full"></div>
                      Handles SSH & HTTPS auth
                    </li>
                    <li className="flex items-center gap-2">
                      <div className="w-1.5 h-1.5 bg-orange-400 rounded-full"></div>
                      One-click web navigation
                    </li>
                  </ul>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}