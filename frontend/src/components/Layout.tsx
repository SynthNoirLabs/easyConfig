import React, { ReactNode } from 'react';
import './Layout.css';

interface LayoutProps {
  sidebar: ReactNode;
  children: ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ sidebar, children }) => {
  return (
    <div className="layout-container">
      <aside className="layout-sidebar">
        {sidebar}
      </aside>
      <main className="layout-main">
        {children}
      </main>
    </div>
  );
};

export default Layout;
