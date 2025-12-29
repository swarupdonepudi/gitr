
import React, { useState, useEffect } from 'react';
import { Button } from "@/components/ui/button";
import { 
  Github, 
  Terminal, 
  ChevronDown
} from 'lucide-react';

import Hero from '../components/home/Hero';
import ProblemStatement from '../components/home/ProblemStatement';
import FeatureGrid from '../components/home/FeatureGrid';
import HowItWorks from '../components/home/HowItWorks';
import Quickstart from '../components/home/Quickstart';
import CLIReference from '../components/home/CLIReference';
import CompareSection from '../components/home/CompareSection';
import FAQ from '../components/home/FAQ';
import Footer from '../components/home/Footer';

export default function Home() {
  const [activeSection, setActiveSection] = useState('hero');

  useEffect(() => {
    const handleScroll = () => {
      const sections = ['hero', 'why', 'how', 'quickstart', 'cli', 'compare', 'faq'];
      const scrollPosition = window.scrollY + 100;

      for (const section of sections) {
        const element = document.getElementById(section);
        if (element) {
          const { offsetTop, offsetHeight } = element;
          if (scrollPosition >= offsetTop && scrollPosition < offsetTop + offsetHeight) {
            setActiveSection(section);
            break;
          }
        }
      }
    };

    window.addEventListener('scroll', handleScroll);
    return () => window.removeEventListener('scroll', handleScroll);
  }, []);

  const scrollToSection = (sectionId) => {
    const element = document.getElementById(sectionId);
    if (element) {
      element.scrollIntoView({ behavior: 'smooth' });
    }
  };

  return (
    <div className="min-h-screen bg-black text-white relative">
      {/* Background Effects */}
      <div className="fixed inset-0 bg-gradient-to-br from-blue-950/20 via-black to-purple-950/20" />
      <div className="fixed inset-0 bg-[radial-gradient(circle_at_center,transparent_0%,rgba(0,0,0,0.8)_100%)]" />
      
      {/* Navigation */}
      <nav className="fixed top-0 w-full z-50 bg-black/80 backdrop-blur-xl border-b border-gray-800/50">
        <div className="max-w-7xl mx-auto px-6 py-4">
          <div className="flex justify-between items-center">
            <div className="flex items-center gap-3">
              <div className="w-8 h-8 bg-gradient-to-r from-blue-500 to-purple-500 rounded-lg flex items-center justify-center">
                <Terminal className="w-5 h-5 text-white" />
              </div>
              <span className="text-xl font-bold bg-gradient-to-r from-blue-400 to-purple-400 bg-clip-text text-transparent">
                gitr
              </span>
            </div>
            
            <div className="hidden md:flex items-center gap-8">
              {[
                { id: 'why', label: 'Why' },
                { id: 'how', label: 'How it works' },
                { id: 'quickstart', label: 'Quickstart' },
                { id: 'cli', label: 'CLI' },
                { id: 'compare', label: 'Compare' },
                { id: 'faq', label: 'FAQ' }
              ].map((item) => (
                <button
                  key={item.id}
                  onClick={() => scrollToSection(item.id)}
                  className={`text-sm transition-colors ${
                    activeSection === item.id 
                      ? 'text-blue-400' 
                      : 'text-gray-400 hover:text-white'
                  }`}
                >
                  {item.label}
                </button>
              ))}
            </div>

            <div className="flex items-center gap-4">
              <Button variant="ghost" size="sm" asChild>
                <a href="https://github.com/swarupdonepudi/gitr" target="_blank" rel="noopener noreferrer">
                  <Github className="w-4 h-4 mr-2" />
                  GitHub
                </a>
              </Button>
            </div>
          </div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="relative pt-20">
        <section id="hero">
          <Hero />
        </section>

        {/* New Problem Statement Section */}
        <section className="py-0">
          <ProblemStatement />
        </section>

        <section id="why" className="py-20">
          <FeatureGrid />
        </section>

        <section id="how" className="py-20">
          <HowItWorks />
        </section>

        <section id="quickstart" className="py-20">
          <Quickstart />
        </section>

        <section id="cli" className="py-20">
          <CLIReference />
        </section>

        <section id="compare" className="py-20">
          <CompareSection />
        </section>

        <section id="faq" className="py-20">
          <FAQ />
        </section>
      </main>

      <Footer />
    </div>
  );
}
