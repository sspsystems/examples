from flask import Flask, request, jsonify
from functools import wraps
from datetime import datetime
import os
from dotenv import load_dotenv

load_dotenv()

app = Flask(__name__)

# Configuration
SSP_API_KEY = os.getenv('SSP_API_KEY')
PORT = int(os.getenv('PORT', 3000))

# Authentication decorator
def require_auth(f):
    @wraps(f)
    def decorated_function(*args, **kwargs):
        api_key = request.headers.get('X-API-Key')

        if not api_key or api_key != SSP_API_KEY:
            return jsonify({
                'error': True,
                'message': 'Unauthorized'
            }), 401

        return f(*args, **kwargs)
    return decorated_function

# Health check - REQUIRED
@app.route('/health', methods=['GET'])
def health():
    return jsonify({
        'status': 'ok',
        'version': '1.0.0',
        'timestamp': datetime.utcnow().isoformat() + 'Z'
    })

# Capabilities - REQUIRED
@app.route('/capabilities', methods=['GET'])
def capabilities():
    return jsonify({
        'supported_methods': ['your_methods_here'],
        'supported_currencies': ['USD', 'INR'],
        'features': ['feature1', 'feature2']
    })

# Example endpoint - Implement your business logic here
@app.route('/your-endpoint', methods=['POST'])
@require_auth
def your_endpoint():
    try:
        data = request.get_json()

        # TODO: Implement your logic here
        # 1. Validate input
        # 2. Process request
        # 3. Return response

        return jsonify({
            'success': True,
            'message': 'Request processed successfully',
            'data': {
                # Your response data
            }
        })

    except Exception as e:
        app.logger.error(f'Error: {str(e)}')
        return jsonify({
            'error': True,
            'message': str(e)
        }), 500

# Error handlers
@app.errorhandler(404)
def not_found(error):
    return jsonify({
        'error': True,
        'message': 'Endpoint not found'
    }), 404

@app.errorhandler(500)
def internal_error(error):
    return jsonify({
        'error': True,
        'message': 'Internal server error'
    }), 500

if __name__ == '__main__':
    print(f'SSP Plugin running on port {PORT}')
    print(f'Health: http://localhost:{PORT}/health')
    app.run(host='0.0.0.0', port=PORT, debug=os.getenv('DEBUG', 'False') == 'True')
