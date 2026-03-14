import { fireEvent, render, screen } from "@testing-library/react";
import { beforeEach, expect, type Mock, test, vi } from "vitest";
import AddConfigModal from "../../../src/components/AddConfigModal";
import { useConfig } from "../../../src/context/ConfigContext";

// Mock hooks and Wails bridge
vi.mock("../../../src/context/ConfigContext");
const mockedUseConfig = vi.mocked(useConfig);

const mockOnClose = vi.fn();
const mockOnSuccess = vi.fn();

beforeEach(() => {
  window.go = {
    main: {
      App: {
        CreateConfig: vi.fn(),
      },
    },
  } as Window["go"];
  mockedUseConfig.mockReset();
});

test("AddConfigModal renders when open", () => {
  mockedUseConfig.mockReturnValue({ configs: [] } as unknown as ReturnType<
    typeof useConfig
  >);
  render(
    <AddConfigModal
      isOpen={true}
      onClose={mockOnClose}
      onSuccess={mockOnSuccess}
    />,
  );
  expect(screen.getByText("Add Configuration")).toBeInTheDocument();
});

test("AddConfigModal handles form submission", async () => {
  mockedUseConfig.mockReturnValue({ configs: [] } as unknown as ReturnType<
    typeof useConfig
  >);
  (window.go?.main?.App?.CreateConfig as Mock).mockResolvedValue(undefined);

  render(
    <AddConfigModal
      isOpen={true}
      onClose={mockOnClose}
      onSuccess={mockOnSuccess}
    />,
  );

  fireEvent.click(screen.getByText("Create"));

  expect(window.go?.main?.App?.CreateConfig).toHaveBeenCalledWith(
    "Claude Code",
    "global",
    "",
  );
});
