import { describe, it, expect } from "vitest";
import { mapAuthDTOToUserViewModel } from "../mappers";

describe("mapAuthDTOToUserViewModel", () => {
  it("should map AuthDTO with user email and name", () => {
    const dto = {
      token: "jwt_token_123",
      user: {
        id: "usr_99",
        email: "alice@cloudstorex.com",
        full_name: "Alice Engineer",
      },
    };
    const vm = mapAuthDTOToUserViewModel(dto);
    expect(vm.id).toBe("usr_99");
    expect(vm.email).toBe("alice@cloudstorex.com");
    expect(vm.fullName).toBe("Alice Engineer");
    expect(vm.avatarUrl).toContain("Alice%20Engineer");
  });

  it("should use fallback email if user object is undefined", () => {
    const dto = { token: "jwt_token_456" };
    const vm = mapAuthDTOToUserViewModel(dto, "dev@cloudstorex.com");
    expect(vm.email).toBe("dev@cloudstorex.com");
    expect(vm.fullName).toBe("dev");
  });
});
