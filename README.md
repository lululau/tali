# tali

A Terminal User Interface (TUI) application for managing Alibaba Cloud resources. Navigate and inspect your cloud infrastructure directly from the command line.

## Features

- **ECS Instances**: View instance details, status, IP addresses, and configurations
- **DNS Management**: Browse AliDNS domains and their DNS records
- **SLB Instances**: Monitor Server Load Balancer instances and their configurations
- **RDS Instances**: Inspect RDS database instances and their properties
- **OSS Management**: Browse OSS buckets and objects (currently stubbed)

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
- `Ctrl+C` - Quit application
- `Esc` or `q` - Go back to previous screen/menu
- `Enter` - Select item or view details

#### List Navigation
- `j` - Move down
- `k` - Move up
- `Enter` - Select item for detailed view

#### Main Menu Options
- `1` - View ECS Instances
- `2` - DNS Management
- `3` - View SLB Instances
- `4` - OSS Management (currently stubbed)
- `5` - View RDS Instances
- `q` - Quit application

### Screens and Features

#### ECS Instances
- Lists all ECS instances with ID, status, IP address, and name
- Select an instance to view detailed information including:
  - Instance specifications
  - Network configuration
  - Security groups
  - Storage details

#### DNS Management
- Browse all domains in your account
- View record count and version information
- Select a domain to view all DNS records
- See record types, values, TTL, and status

#### SLB Instances
- List all Server Load Balancer instances
- View SLB ID, name, IP address, type, and status
- Select for detailed configuration view

#### RDS Instances
- Browse all RDS database instances
- View engine type, version, instance class, and status
- Select for detailed database configuration including:
  - Connection strings and ports
  - Storage information
  - Network configuration
  - Backup and maintenance windows

## Required Permissions

Your Alibaba Cloud Access Key needs the following permissions:

- **ECS**: `ecs:DescribeInstances`
- **DNS**: `alidns:DescribeDomains`, `alidns:DescribeDomainRecords`
- **SLB**: `slb:DescribeLoadBalancers`
- **RDS**: `rds:DescribeDBInstances`
- **OSS**: `oss:ListBuckets`, `oss:ListObjects` (when implemented)

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