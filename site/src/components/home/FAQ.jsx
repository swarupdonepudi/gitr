import React, { useState } from 'react';
import { Card, CardContent } from "@/components/ui/card";
import { ChevronDown, ChevronRight } from 'lucide-react';

const faqs = [
  {
    question: 'Do I need git installed?',
    answer: 'For SSH clones, gitr shells out to git(1). For HTTPS clones, gitr uses go-git and does not require the git binary.'
  },
  {
    question: 'Which providers are supported today?',
    answer: 'github.com, gitlab.com, bitbucket.org are built-in. Custom hosts can be configured via ~/.gitr.yaml.'
  },
  {
    question: 'How are tokens/keys discovered?',
    answer: 'HTTPS tokens: $HOME/.personal_access_tokens/<hostname> or --token flag. SSH: $HOME/.ssh/config (IdentityFile) and Host mappings.'
  },
  {
    question: 'Where is the clone path rooted?',
    answer: 'Default is ~/scm/ but can be configured in ~/.gitr.yaml. Directory creation mirrors host/owner/repo structure.'
  },
  {
    question: 'How does preview work?',
    answer: '--dry prints exactly what gitr will do (provider, host, repo path, URLs) without cloning or opening a browser.'
  },
  {
    question: 'Will it overwrite existing repos?',
    answer: 'No. If the target path exists and contains .git, gitr skips cloning to prevent data loss.'
  },
  {
    question: 'What license is gitr under?',
    answer: 'Apache-2.0. Free for commercial and personal use.'
  }
];

export default function FAQ() {
  const [openIndex, setOpenIndex] = useState(null);

  const toggleFAQ = (index) => {
    setOpenIndex(openIndex === index ? null : index);
  };

  return (
    <div className="max-w-4xl mx-auto px-6">
      <div className="text-center mb-16">
        <h2 className="text-4xl md:text-5xl font-bold mb-6 bg-gradient-to-r from-white to-gray-300 bg-clip-text text-transparent">
          Frequently asked questions
        </h2>
        <p className="text-xl text-gray-400">
          Everything you need to know about gitr
        </p>
      </div>

      <div className="space-y-4">
        {faqs.map((faq, index) => (
          <Card 
            key={index} 
            className="bg-gray-900/50 border-gray-700 hover:border-gray-600 transition-colors cursor-pointer"
            onClick={() => toggleFAQ(index)}
          >
            <CardContent className="p-6">
              <div className="flex items-center justify-between">
                <h3 className="text-lg font-semibold text-white pr-4">
                  {faq.question}
                </h3>
                {openIndex === index ? (
                  <ChevronDown className="w-5 h-5 text-gray-400 flex-shrink-0" />
                ) : (
                  <ChevronRight className="w-5 h-5 text-gray-400 flex-shrink-0" />
                )}
              </div>
              
              {openIndex === index && (
                <div className="mt-4 pt-4 border-t border-gray-700">
                  <p className="text-gray-300 leading-relaxed">
                    {faq.answer}
                  </p>
                </div>
              )}
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Contact Section */}
      <div className="mt-16 text-center bg-gradient-to-r from-gray-900/50 to-gray-800/50 backdrop-blur-xl border border-gray-700 rounded-3xl p-8">
        <h3 className="text-2xl font-bold text-white mb-4">
          Still have questions?
        </h3>
        <p className="text-gray-400 mb-6">
          Check out the documentation or report an issue
        </p>
        <div className="flex flex-col sm:flex-row gap-4 justify-center">
          <a 
            href="https://github.com/swarupdonepudi/gitr" 
            target="_blank" 
            rel="noopener noreferrer"
            className="px-6 py-3 bg-gray-800 hover:bg-gray-700 text-white rounded-lg transition-colors"
          >
            View Documentation
          </a>
          <a 
            href="https://github.com/swarupdonepudi/gitr/issues" 
            target="_blank" 
            rel="noopener noreferrer"
            className="px-6 py-3 border border-gray-600 hover:bg-gray-800 text-gray-300 rounded-lg transition-colors"
          >
            Report an Issue
          </a>
        </div>
      </div>
    </div>
  );
}