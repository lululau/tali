# tali

A Terminal User Interface (TUI) application for managing Alibaba Cloud resources. Navigate and inspect your cloud infrastructure directly from the command line.

## Features

- **ECS Instances**: View instance details with zone, CPU/RAM configuration, private/public IPs, and full JSON details
- **DNS Management**: Browse AliDNS domains and their DNS records
- **SLB Instances**: Monitor Server Load Balancer instances and their configurations with JSON details
- **RDS Instances**: Inspect RDS database instances and their properties with JSON details
- **RDS Database Management**: 
  - View databases and accounts for each RDS instance
  - Press `D` on RDS instance to view databases
  - Press `A` on RDS instance to view accounts
  - Full JSON details for databases and accounts
- **OSS Management**: Browse OSS buckets and objects with pagination and full details
- **Interactive Features**: 
  - Copy any data as JSON to clipboard with double-y
  - Edit JSON data in external nvim editor
  - Mouse text selection in detail views
  - Real-time pagination info in mode line

## Prerequisites

- Go 1.22.2 or later
- Valid Alibaba Cloud account with API access
- Access Key ID and Secret with appropriate permissions for the services you want to manage

## Installation

### From Source

1. Clone the repository:
```bash
git clone <repository-url>
cd tali
```

2. Build the application:
```bash
go build -o tali cmd/main.go
```

3. (Optional) Install globally:
```bash
go install ./cmd
```

## Configuration

### Environment Variables

The application requires the following environment variables to be set:

```bash
export ALIBABA_CLOUD_ACCESS_KEY_ID="your-access-key-id"
export ALIBABA_CLOUD_ACCESS_KEY_SECRET="your-access-key-secret"
export ALIBABA_CLOUD_REGION_ID="your-region-id"
```

Optional:
```bash
export ALIBABA_CLOUD_OSS_ENDPOINT="your-oss-endpoint"
```

### Setting up Environment Variables

Create a `.env` file or add to your shell profile:

```bash
# Example for ~/.bashrc or ~/.zshrc
export ALIBABA_CLOUD_ACCESS_KEY_ID="LTAI4..."
export ALIBABA_CLOUD_ACCESS_KEY_SECRET="ABC123..."
export ALIBABA_CLOUD_REGION_ID="cn-hangzhou"
export ALIBABA_CLOUD_OSS_ENDPOINT="oss-cn-hangzhou.aliyuncs.com"
```

Common region IDs:
- `cn-hangzhou` - China (Hangzhou)
- `cn-shanghai` - China (Shanghai)
- `cn-beijing` - China (Beijing)
- `us-west-1` - US West (Silicon Valley)
- `ap-southeast-1` - Asia Pacific (Singapore)

## Usage

### Running the Application

```bash
./tali
```

Or if installed globally:
```bash
tali
```

### Navigation

The application uses vim-style keyboard navigation:

#### Global Controls
- `Ctrl+C` or `Q` - Quit application
- `Esc` or `q` - Go back to previous screen/menu
- `Enter` - Select item or view details
- `O` - Open profile selection dialog

#### Profile Management
- The mode line at the bottom shows the current active profile
- Press `O` to open the profile selection dialog
- In the profile dialog:
  - Use `j/k` to navigate up/down
  - Press `Enter` to select a profile
  - Press `q` or `Esc` to cancel
- After switching profiles:
  - All client connections are automatically recreated
  - All cached data is cleared
  - Returns to main menu
  - New credentials take effect immediately

#### OSS Object Pagination
- `[` - Previous page
- `]` - Next page  
- `0` - Go to first page
- Page information displayed in the bottom-right mode line

#### Data Copying and Editing
- `yy` (double-y) - Copy current row JSON data to clipboard (in tables)
- `yy` (double-y) - Copy complete JSON data to clipboard (in detail views)
- `e` - Open JSON data in nvim for editing (in detail views)
- Detail views support mouse text selection

#### Search Functionality
- `/` - Enter search mode (shows search bar at bottom of screen, vim-style)
- `Enter` - Execute search and highlight matches
- `Escape` - Exit search mode
- `n` - Go to next search result
- `p` - Go to previous search result
- Search is case-insensitive by default
- Works in all table views and JSON detail views
- Search results are highlighted in yellow
- Search interface mimics vim's behavior with bottom search bar

#### List Navigation
- `j` - Move down
- `k` - Move up
- `Enter` - Select item for detailed view

#### RDS Instance Management
- `D` - View databases for selected RDS instance
- `A` - View accounts for selected RDS instance

#### Main Menu Options
- `1` - View ECS Instances
- `2` - View Security Groups
- `3` - DNS Management
- `4` - View SLB Instances
- `5` - OSS Management
- `6` - View RDS Instances
- `Q` - Quit application

### Screens and Features

#### ECS Instances
- Lists all ECS instances with ID, status, zone, CPU/RAM configuration, private IP, public IP, and name
- Select an instance to view complete JSON details including:
  - Instance specifications
  - Network configuration
  - Security groups
  - Storage details
  - All available metadata

#### Security Groups
- Lists all ECS security groups with ID, name, description, VPC ID, type, and creation time
- Select a security group to view complete JSON details including:
  - Security group rules (ingress and egress)
  - Associated instances
  - Network configuration
  - All available metadata

#### DNS Management
- Browse all domains in your account
- View record count and version information
- Select a domain to view all DNS records
- See record types, values, TTL, and status

#### SLB Instances
- List all Server Load Balancer instances
- View SLB ID, name, IP address, type, and status
- Select for complete JSON configuration view

#### RDS Instances
- Browse all RDS database instances
- View engine type, version, instance class, and status
- Select for complete JSON database configuration including:
  - Connection strings and ports
  - Storage information
  - Network configuration
  - Backup and maintenance windows
  - All available metadata

#### RDS Database Management
- View databases and accounts for each RDS instance
- Press `D` on RDS instance to view databases
- Press `A` on RDS instance to view accounts
- Full JSON details for databases and accounts

#### OSS Management
- Browse all OSS buckets with name, location, creation date, and storage class
- Select a bucket to view all objects with key, size, last modified date, storage class, and ETag
- Select an object to view complete JSON metadata

## Required Permissions

Your Alibaba Cloud Access Key needs the following permissions:

- **ECS**: `ecs:DescribeInstances`, `ecs:DescribeSecurityGroups`
- **DNS**: `alidns:DescribeDomains`, `alidns:DescribeDomainRecords`
- **SLB**: `slb:DescribeLoadBalancers`
- **RDS**: `rds:DescribeDBInstances`
- **OSS**: `oss:ListBuckets`, `oss:ListObjects`

## Troubleshooting

### Common Issues

1. **Authentication Error**
   - Verify your Access Key ID and Secret are correct
   - Ensure the keys have the required permissions
   - Check that the region ID is valid

2. **No Resources Found**
   - Verify you're using the correct region ID
   - Ensure resources exist in the specified region
   - Check that your account has access to the resources

3. **Network Issues**
   - Ensure you have internet connectivity
   - Check if your firewall allows HTTPS traffic
   - Verify the OSS endpoint (if using OSS features)

### Debug Mode

For additional debugging information, you can check the application logs or run with verbose output by setting:

```bash
export ALIBABA_CLOUD_DEBUG=true
```

## Development

### Dependencies

The project uses Go modules. Key dependencies include:

- `github.com/aliyun/alibaba-cloud-sdk-go` - Alibaba Cloud SDK
- `github.com/rivo/tview` - Terminal UI framework
- `github.com/gdamore/tcell/v2` - Terminal cell manipulation

### Building from Source

```bash
# Download dependencies
go mod download

# Run tests (if any)
go test ./...

# Build
go build -o tali cmd/main.go
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Test thoroughly
5. Submit a pull request

## License

[Add your license information here]

## Support

[Add support information here]