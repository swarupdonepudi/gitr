import React, { useState } from 'react';
import { Button } from "@/components/ui/button";
import { Tabs, TabsContent, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Copy, Check, Terminal, Apple, Monitor } from 'lucide-react';

const CodeBlock = ({ children, filename, language = 'bash' }) => {
  const [copied, setCopied] = useState(false);

  const copyCode = () => {
    const code = children.toString().replace(/^\$ /gm, '');
    navigator.clipboard.writeText(code);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="bg-black/50 border border-gray-800 rounded-xl overflow-hidden">
      {filename && (
        <div className="flex items-center justify-between px-4 py-2 bg-gray-900 border-b border-gray-800">
          <div className="flex items-center gap-2">
            <div className="w-3 h-3 bg-red-500 rounded-full"></div>
            <div className="w-3 h-3 bg-yellow-500 rounded-full"></div>
            <div className="w-3 h-3 bg-green-500 rounded-full"></div>
            <span className="text-gray-400 text-sm ml-2">{filename}</span>
          </div>
          <Button
            variant="ghost"
            size="sm"
            onClick={copyCode}
            className="text-gray-400 hover:text-white h-6 px-2"
          >
            {copied ? (
              <Check className="w-3 h-3 text-green-400" />
            ) : (
              <Copy className="w-3 h-3" />
            )}
          </Button>
        </div>
      )}
      <pre className="p-4 overflow-x-auto">
        <code className="text-gray-300 font-mono text-sm">
          {children}
        </code>
      </pre>
    </div>
  );
};

export default function Quickstart() {
  return (
    <div className="max-w-7xl mx-auto px-6">
      <div className="text-center mb-16">
        <h2 className="text-4xl md:text-5xl font-bold mb-6 bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
          Get started in minutes
        </h2>
        <p className="text-xl text-gray-400 max-w-3xl mx-auto">
          Install gitr and transform your git workflow with four simple steps
        </p>
      </div>

      <div className="max-w-4xl mx-auto">
        {/* Installation */}
        <div className="mb-12">
          <h3 className="text-2xl font-bold text-white mb-6 flex items-center gap-3">
            <span className="w-8 h-8 bg-gradient-to-r from-blue-600 to-purple-600 rounded-lg flex items-center justify-center text-sm font-bold">
              1
            </span>
            Install gitr
          </h3>
          
          <Tabs defaultValue="homebrew" className="w-full">
            <TabsList className="bg-gray-900 border border-gray-700">
              <TabsTrigger value="homebrew" className="flex items-center gap-2">
                <Apple className="w-4 h-4" />
                macOS (Homebrew)
              </TabsTrigger>
              <TabsTrigger value="go" className="flex items-center gap-2">
                <Terminal className="w-4 h-4" />
                Go Install
              </TabsTrigger>
              <TabsTrigger value="binary" className="flex items-center gap-2">
                <Monitor className="w-4 h-4" />
                Binary
              </TabsTrigger>
            </TabsList>
            
            <TabsContent value="homebrew" className="mt-4">
              <CodeBlock filename="Terminal">
{`# macOS (Homebrew)
$ brew install swarupdonepudi/tap/gitr`}
              </CodeBlock>
            </TabsContent>
            
            <TabsContent value="go" className="mt-4">
              <CodeBlock filename="Terminal">
{`# Go Install (any platform)
$ go install github.com/swarupdonepudi/gitr@latest`}
              </CodeBlock>
            </TabsContent>
            
            <TabsContent value="binary" className="mt-4">
              <CodeBlock filename="Terminal">
{`# Download from releases
$ curl -L https://github.com/swarupdonepudi/gitr/releases/latest/download/gitr-$(uname -s)-$(uname -m) -o gitr
$ chmod +x gitr
$ sudo mv gitr /usr/local/bin/`}
              </CodeBlock>
            </TabsContent>
          </Tabs>
        </div>

        {/* Clone Example */}
        <div className="mb-12">
          <h3 className="text-2xl font-bold text-white mb-6 flex items-center gap-3">
            <span className="w-8 h-8 bg-gradient-to-r from-green-600 to-emerald-600 rounded-lg flex items-center justify-center text-sm font-bold">
              2
            </span>
            Clone to a deterministic path
          </h3>
          
          <CodeBlock filename="Terminal">
{`$ gitr clone https://github.com/owner/repo
# Clones to ~/scm/github.com/owner/repo
# Copies: cd ~/scm/github.com/owner/repo`}
          </CodeBlock>
          
          <div className="mt-4 p-4 bg-blue-950/30 border border-blue-800/50 rounded-xl">
            <p className="text-blue-200 text-sm">
              💡 <strong>Pro tip:</strong> The exact path is copied to your clipboard automatically. 
              Just paste with Cmd+V to navigate there instantly.
            </p>
          </div>
        </div>

        {/* Web Navigation */}
        <div className="mb-12">
          <h3 className="text-2xl font-bold text-white mb-6 flex items-center gap-3">
            <span className="w-8 h-8 bg-gradient-to-r from-purple-600 to-pink-600 rounded-lg flex items-center justify-center text-sm font-bold">
              3
            </span>
            Jump to repo web pages
          </h3>
          
          <CodeBlock filename="Terminal - Run inside a repo">
{`$ gitr web           # repo home
$ gitr prs           # PRs/MRs
$ gitr pipelines     # GitHub Actions / GitLab Pipelines / Bitbucket Pipelines
$ gitr issues        # issues
$ gitr rem           # current branch in the web UI
$ gitr commits       # commits of current branch`}
          </CodeBlock>
        </div>

        {/* Preview Mode */}
        <div className="mb-12">
          <h3 className="text-2xl font-bold text-white mb-6 flex items-center gap-3">
            <span className="w-8 h-8 bg-gradient-to-r from-orange-600 to-red-600 rounded-lg flex items-center justify-center text-sm font-bold">
              4
            </span>
            Preview without changing anything
          </h3>
          
          <CodeBlock filename="Terminal">
{`$ gitr --dry clone https://github.com/owner/repo
$ gitr --dry web`}
          </CodeBlock>
          
          <div className="mt-4 p-4 bg-gray-900/50 border border-gray-700 rounded-xl">
            <p className="text-gray-300 text-sm">
              <strong>Authentication:</strong> HTTPS tokens can be passed with --token or read from 
              <code className="bg-gray-800 px-1 rounded text-xs ml-1">$HOME/.personal_access_tokens/&lt;hostname&gt;</code>. 
              SSH keys and aliases are discovered from <code className="bg-gray-800 px-1 rounded text-xs">$HOME/.ssh/config</code>.
            </p>
          </div>
        </div>

        {/* Next Steps */}
        <div className="bg-gradient-to-r from-gray-900/50 to-gray-800/50 backdrop-blur-xl border border-gray-700 rounded-3xl p-8 text-center">
          <h3 className="text-2xl font-bold text-white mb-4">Ready to level up your workflow?</h3>
          <p className="text-gray-400 mb-6">
            Explore advanced configuration and CLI commands
          </p>
          <div className="flex flex-col sm:flex-row gap-4 justify-center">
            <Button 
              size="lg"
              className="bg-gradient-to-r from-blue-600 to-purple-600 hover:from-blue-700 hover:to-purple-700"
              onClick={() => document.getElementById('cli').scrollIntoView({ behavior: 'smooth' })}
            >
              View CLI Reference
            </Button>
            <Button 
              variant="outline" 
              size="lg"
              className="border-gray-600 text-gray-300 hover:bg-gray-800"
              asChild
            >
              <a href="https://github.com/swarupdonepudi/gitr" target="_blank" rel="noopener noreferrer">
                Read Full Docs
              </a>
            </Button>
          </div>
        </div>
      </div>
    </div>
  );
}