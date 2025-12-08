# Configuration Wizards

`easyConfig` now includes an interactive wizard to help you create and update configurations for your favorite tools.

## How to Use the Wizard

1.  **Launch `easyConfig`**: Open the application to see the list of supported providers.
2.  **Find a Provider**: Locate the provider you want to configure (e.g., "Jules").
3.  **Click "Wizard"**: If a provider has a wizard, a "Wizard" button will appear next to its name. Click it to start.
4.  **Follow the Prompts**: The wizard will guide you through a series of steps to create your configuration file.
5.  **Finish**: Once you've completed the steps, the wizard will save the new configuration file.

## Example: Configuring Jules

The Jules provider includes a wizard that helps you create your `data.json` file. The wizard will ask for your name and then create the file with the appropriate content.

## For Developers: Adding a Wizard to a Provider

To add a wizard to your own provider, you need to implement the `Wizard` interface in your provider's Go code.

1.  **Implement the `Wizard` Interface**: Create a struct that implements the `Start`, `Next`, and `Cancel` methods.
2.  **Return the Wizard**: In your provider's `GetWizard` method, return an instance of your new wizard struct.

For a complete example, see the `JulesWizard` implementation in `pkg/config/provider_jules.go`.
