# Contributing to ServerCommander

Thank you for your interest in contributing to **ServerCommander**!  
We welcome contributions of all kinds, including **bug reports, feature requests, documentation improvements, and code contributions**.

---

## How to Contribute

### 1. Fork the Repository

Start by **forking** this repository to your GitHub account:

- Navigate to the [ServerCommander repository](https://github.com/VectoDE/ServerCommander).
- Click the **Fork** button in the top right corner.
- Clone the forked repository to your local machine:

    ```bash
    git clone https://github.com/your-username/ServerCommander.git
    cd ServerCommander
    ```

### 2. Create a New Branch

Before making any changes, create a new branch:
    ```bash
    git checkout -b feature-branch-name
    ```

Use a descriptive branch name, such as:

- ```fix-ssh-connection-issue```
- ```add-sftp-file-transfer```
- ```update-documentation```

### 3. Implement Your Changes

Make your changes while following these guidelines:

- **Code Formatting**: Ensure your code follows Go best practices.
- **Commenting**: Add necessary comments for better code readability.
- **Error Handling**: Use structured error handling for robustness.
- **Logging**: Use the internal logging utility for debugging.

### 4. Test Your Changes

Before submitting your changes, ensure everything works properly:
    ```bash
    go test ./...
    ```

Add new test cases in the ```tests/``` directory if applicable.

### 5. Commit Your Changes

Follow a consistent commit message format:
    ```bash
    git commit -m "fix: resolved SSH connection timeout issue"
    ```

Use the following prefixes:

- ```fix```: for bug fixes
- ```feat```: for new features
- ```docs```: for documentation updates
- ```test```: for test improvements
- ```refactor```: for code restructuring

### 6. Push & Create a Pull Request

Push your branch to GitHub:
    ```bash
    git push origin feature-branch-name
    ```

Then, create a **Pull Request (PR)**:

  1. Go to the [ServerCommander repository](https://github.com/VectoDE/ServerCommander).
  2. Click **New Pull Request**.
  3. Select your branch and provide a clear description of your changes.
  4. Submit the PR for review.

## Contribution Guidelines

- Ensure that your changes **do not break existing functionality**.
- Follow the **Go coding standards**.
- Contributions should be **platform-independent** (Windows, Linux, macOS).
- Avoid unnecessary dependencies to keep the project lightweight.
- Be **respectful** and follow community guidelines when discussing issues.

## Reporting Issues

If you encounter a bug or have a feature request, open an issue:

  1. Go to the **Issues** tab in the GitHub repository.
  2. Click **New Issue**.
  3. Provide a **clear and detailed description**:
    - Steps to reproduce (if a bug)
    - Expected vs. actual behavior
    - Possible solution ideas (if applicable)

## Thank You```!```

Your contributions make ServerCommander better! 🚀

For any questions, feel free to join the discussion or contact the maintainers.

Happy coding! 🎉
