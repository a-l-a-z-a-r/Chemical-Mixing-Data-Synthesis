from tensorflow.keras.models import load_model
from tensorflow.keras.utils import plot_model

# Load the model
model = load_model("model.h5")

# Plot and save model architecture to file
plot_model(
    model, 
    to_file="model_architecture.png", 
    show_shapes=True, 
    show_layer_names=True,
    expand_nested=True,
    dpi=100
)

print("Model architecture saved as model_architecture.png")
